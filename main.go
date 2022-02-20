package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"unicode/utf16"

	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/app/command"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account"
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

func main() {
	fmt.Println("Started")

	psqlConn, err := adapter.NewPostgresConnPool()
	PanicIfErr(err)

	accountModule := account.NewModule(account.AccountModuleConfig{
		PsqlConn: psqlConn,
		BasicAuthCredentials: struct {
			User string
			Pass string
		}{
			User: os.Getenv("BASIC_AUTH_USER"),
			Pass: os.Getenv("BASIC_AUTH_PASS"),
		},
	})

	rootRouter := chi.NewRouter()
	accountRouter := chi.NewRouter()
	accountModule.BindRouter(accountRouter)
	rootRouter.Mount("/account/api/v1", accountRouter)

	botToken := os.Getenv(BOT_TOKEN)
	fmt.Printf("botToken = %#v\n", botToken)
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	port := os.Getenv(PORT)

	channelRepo := adapter.NewPostgresChannelRepository(psqlConn)
	messageRepository := adapter.NewPostgresMessageRepository(psqlConn)
	messageSender := adapter.NewMessageSender(bot)
	replyToChannelsHandler = command.NewForwardToChannelsHandler(channelRepo, messageRepository, messageSender)

	updates := getWebhookUpdateChan(rootRouter)

	go http.ListenAndServe(fmt.Sprintf(":%s", port), rootRouter)

	tgUpdateHandler := telegram.NewTelegramBotUpdateHandler(updates, replyToChannelsHandler)
	tgUpdateHandler.Run()
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getWebhookUpdateChan(router *chi.Mux) tgbotapi.UpdatesChannel {
	_, err := bot.Request(tgbotapi.RemoveWebhookConfig{})
	PanicIfErr(err)

	_, err = bot.Request(tgbotapi.NewWebhook("http://hsequeda.com:8484/"))
	PanicIfErr(err)

	ch := make(chan tgbotapi.Update, bot.Buffer)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)

		var update tgbotapi.Update
		json.Unmarshal(bytes, &update)

		ch <- update
	})

	return ch
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
