package main

import (
	"fmt"
	"net/http"
	"os"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/app/command"
	"github.com/stdevHsequeda/SendToChannelsBot/port/telegram"
)

const (
	PORT            = "PORT"
	BOT_TOKEN       = "BOT_TOKEN"
	WEBHOOK_PATH    = "WEBHOOK_PATH"
	HashtagType     = "hashtag"
	TextMentionType = "text_mention"
	MentionType     = "mention"
	CommandType     = "command"
	ADMIN_ID        = "ADMIN_ID"
)

var bot *tgbotapi.BotAPI

var replyToChannelsHandler command.ForwardToChannelsHandler

func init() {
	var err error

	botToken := os.Getenv(BOT_TOKEN)
	fmt.Printf("botToken = %#v\n", botToken)
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	port := os.Getenv(PORT)
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func main() {
	fmt.Println("Started")

	dbConn, err := adapter.NewPostgresConnPool()
	PanicIfErr(err)

	channelRepo := adapter.NewPostgresChannelRepository(dbConn)
	messageRepository := adapter.NewPostgresMessageRepository(dbConn)
	messageSender := adapter.NewMessageSender(bot)
	replyToChannelsHandler = command.NewForwardToChannelsHandler(channelRepo, messageRepository, messageSender)

	updates, err := getUpdateCh()
	if err != nil {
		panic(err)
	}

	tgUpdateHandler := telegram.NewTelegramBotUpdateHandler(updates, replyToChannelsHandler)

	tgUpdateHandler.Run()
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getUpdateCh() (tgbotapi.UpdatesChannel, error) {
	webhookPath := os.Getenv(WEBHOOK_PATH)
	if webhookPath != "" {
		_, err := bot.Request(tgbotapi.NewWebhook(webhookPath))
		if err != nil {
			panic(err)
		}

		return bot.ListenForWebhook("/"), nil
	} else {
		return bot.GetUpdatesChan(tgbotapi.UpdateConfig{
			Offset:  0,
			Timeout: 0,
		}), nil
	}
}

func containsHashtag(entities []tgbotapi.MessageEntity) bool {
	for _, entity := range entities {
		if entity.Type == HashtagType {
			return true
		}
	}

	return false
}

func getHashtagList(text string, entities []tgbotapi.MessageEntity) []string {
	var textUtf16 = utf16.Encode([]rune(text))
	hashtagList := make([]string, 0)
	for _, entity := range entities {
		hashtag := string(utf16.Decode(textUtf16[entity.Offset : entity.Offset+entity.Length]))
		hashtagList = append(hashtagList, hashtag)
	}

	return hashtagList
}

func getRefLink(user *tgbotapi.User) string {
	var name string
	var link string
	name = user.UserName
	link = fmt.Sprintf("http://t.me/%s", user.UserName)
	if user.UserName == "" {
		name = user.FirstName
		link = fmt.Sprintf("tg://user?id=%d", user.ID)
	}
	return fmt.Sprintf("\n <a href=\"%s\">Escribir al autor👤(%s)</a> ", link, name)
}

func getUserMention(entities []tgbotapi.MessageEntity, text []rune) (*tgbotapi.User, bool) {
	var textUtf16 = utf16.Encode(text)
	for _, entity := range entities {
		if entity.Type == TextMentionType {
			return entity.User, true
		}
		if entity.Type == MentionType {
			return &tgbotapi.User{
				UserName: string(utf16.Decode(textUtf16[entity.Offset+1 : entity.Offset+entity.Length])),
			}, true
		}
	}
	return nil, false
}
