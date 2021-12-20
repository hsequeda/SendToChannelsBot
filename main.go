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

// var db *PostgresDatabase

// var hashtags map[string][]int64
// var messages map[string][]ChannelMessage

// func getCapasity() int {
// 	// size := 0
// 	// for hashtag, chIdList := range hashtags {
// 	// 	size += int(unsafe.Sizeof(hashtag))
// 	// 	for _, chId := range chIdList {
// 	// 		size += int(unsafe.Sizeof(chId))
// 	// 	}
// 	// }

// 	// for messageId, channelMessageList := range messages {
// 	// 	size += int(unsafe.Sizeof(messageId))
// 	// 	for _, chMsgId := range channelMessageList {
// 	// 		size += int(unsafe.Sizeof(chMsgId))
// 	// 	}
// 	// }
// 	// return size
// }

var replyToChannelsHandler command.ForwardToChannelsHandler

func init() {
	var err error
	// db, err = InitDb()
	// if err != nil {
	// 	panic(err)
	// }

	// hashtags, err = db.ListHashtags()
	// if err != nil {
	// 	panic(err)
	// }

	// messages, err = db.ListMessages()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%#v bits\n", getCapasity())

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

	// channelRepo := adapter.NewMemoryChannelRepository()
	channelRepo := adapter.NewPostgresChannelRepository(memRepo)
	messageSender := adapter.NewMessageSender(bot)
	replyToChannelsHandler = command.NewForwardToChannelsHandler(channelRepo, messageSender)

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

// func resolveUpdate(update *tgbotapi.Update) {
// 	switch {
// 	case update.Message != nil:
// 		// sendMes(update.Message.Chat.ID)
// 		resolveMessage(update.Message)
// 	case update.EditedMessage != nil:
// 		// resolveEditMessage(update.EditedMessage)
// 	case update.ChannelPost != nil:
// 		// resolveChannelPost(update.ChannelPost)
// 	}
// }

// func resolveMessage(message *tgbotapi.Message) {
// 	if message.IsCommand() {
// 		if err := handleCommand(message); err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		return
// 	}

// 	if message.Entities != nil || message.CaptionEntities != nil {
// 		channelMessageList, err := resolveHashtagType(message)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}

// 		if channelMessageList != nil {
// 			messages, err = db.UpdateMessages(strconv.FormatInt(int64(message.MessageID), 10), channelMessageList)
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 		}
// 	}
// }

// func resolveHashtagType(message *tgbotapi.Message) ([]ChannelMessage, error) {
// 	if message.Text != "" {
// 		return sendMessageFromText(message)
// 	}

// 	if message.Caption != "" {
// 		return sendMessageFromCaption(message)
// 	}

// 	return nil, nil
// }

// func getAllData() string {
// 	result := map[int64][]string{}
// 	for key, value := range hashtags {
// 		for _, value2 := range value {
// 			result[value2] = append(result[value2], key)
// 		}
// 	}

// 	var result2 string
// 	for key, value := range result {
// 		chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
// 			ChatConfig: tgbotapi.ChatConfig{ChatID: key},
// 		})

// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}

// 		result2 += fmt.Sprintf("%s-> %s \n", chat.UserName, strings.Join(value, ", "))
// 	}

// 	return result2
// }

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

// func getChannelList(hashtagList []string) []int64 {
// 	channelList := make([]int64, 0)
// 	channelFilter := make(map[int64]struct{})
// 	for _, hashtag := range hashtagList {
// 		// Get list of channels subscribed to a hashtag
// 		if subsChannels, ok := hashtags[hashtag]; ok {
// 			for _, ch := range subsChannels {
// 				if _, ok := channelFilter[ch]; !ok {
// 					channelList = append(channelList, ch)
// 					channelFilter[ch] = struct{}{}
// 				}
// 			}
// 		}
// 	}

// 	return channelList
// }

// func handleCommand(message *tgbotapi.Message) error {
// 	switch {
// 	case regexp.MustCompile(`(?m)^\/list( \w+|$)`).MatchString(message.Text):
// 		adminId := os.Getenv(ADMIN_ID)
// 		if strconv.FormatInt(message.Chat.ID, 10) == adminId {
// 			_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, getAllData()))
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func sendMessageFromText(message *tgbotapi.Message) ([]ChannelMessage, error) {
// 	hashtagList := getHashtagList(message.Text, message.Entities)

// 	// channelIdList := getChannelList(hashtagList)
// 	// var user *tgbotapi.User
// 	// user = message.From
// 	// if userMention, ok := getUserMention(message.Entities, []rune(message.Text)); ok {
// 	// 	user = userMention
// 	// }

// 	err := replyToChannelsHandler.Handle(context.TODO(), command.ForwardToChannels{
// 		Text:        message.Text,
// 		HashtagList: hashtagList,
// 		UserName:    message.From.UserName,
// 		// UserId:      strconv.Itoa(message.From.ID),
// 		// MessageId:   message.MessageID,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return []ChannelMessage{}, nil
// 	// var channelMessageList = make([]ChannelMessage, 0)
// 	// for _, chId := range channelIdList {
// 	// 	toSend := tgbotapi.NewMessage(chId, fmt.Sprintf("%s\n%s", message.Text, getRefLink(user)))
// 	// 	toSend.ParseMode = "html"
// 	// 	toSend.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
// 	// 		tgbotapi.NewInlineKeyboardRow(
// 	// 			tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
// 	// 		),
// 	// 	)

// 	// 	msg, err := bot.Send(toSend)
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}

// 	// 	channelMessageList = append(channelMessageList, ChannelMessage{ChannelId: chId, MessageId: int64(msg.MessageID)})
// 	// }
// }

// func sendMessageFromCaption(message *tgbotapi.Message) ([]ChannelMessage, error) {
// 	hashtagList := getHashtagList(message.Caption, message.CaptionEntities)
// 	channelIdList := getChannelList(hashtagList)
// 	var user *tgbotapi.User
// 	user = message.From
// 	if userMention, ok := getUserMention(message.Entities, []rune(message.Caption)); ok {
// 		user = userMention
// 	}
// 	refLink := getRefLink(user)
// 	var channelMessageList = make([]ChannelMessage, 0)
// 	for _, chId := range channelIdList {
// 		var toSend tgbotapi.Chattable
// 		if message.Photo != nil {
// 			photoConfig := tgbotapi.NewPhotoShare(chId, message.Photo[0].FileID)
// 			photoConfig.Caption = fmt.Sprintf("%s\n%s", message.Caption, refLink)
// 			photoConfig.ParseMode = "html"
// 			photoConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
// 				tgbotapi.NewInlineKeyboardRow(
// 					tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
// 				),
// 			)
// 			toSend = photoConfig
// 		}

// 		if message.Audio != nil {
// 			audioConfig := tgbotapi.NewAudioShare(chId, message.Audio.FileID)
// 			audioConfig.Caption = fmt.Sprintf("%s\n%s", message.Caption, refLink)
// 			audioConfig.ParseMode = "html"
// 			audioConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
// 				tgbotapi.NewInlineKeyboardRow(
// 					tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
// 				),
// 			)
// 			toSend = audioConfig
// 		}

// 		if message.Video != nil {
// 			videoConfig := tgbotapi.NewVideoShare(chId, message.Video.FileID)
// 			videoConfig.Caption = fmt.Sprintf("%s\n%s", message.Caption, refLink)
// 			videoConfig.ParseMode = "html"
// 			videoConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
// 				tgbotapi.NewInlineKeyboardRow(
// 					tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
// 				),
// 			)
// 			toSend = videoConfig
// 		}

// 		if message.Document != nil {
// 			documentConfig := tgbotapi.NewDocumentShare(chId, message.Document.FileID)
// 			documentConfig.Caption = fmt.Sprintf("%s\n%s", message.Caption, refLink)
// 			documentConfig.ParseMode = "html"
// 			documentConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
// 				tgbotapi.NewInlineKeyboardRow(
// 					tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
// 				),
// 			)
// 			toSend = documentConfig
// 		}

// 		msg, err := bot.Send(toSend)
// 		if err != nil {
// 			return nil, err
// 		}

// 		channelMessageList = append(channelMessageList, ChannelMessage{ChannelId: chId, MessageId: int64(msg.MessageID)})
// 	}
// 	return channelMessageList, nil
// }

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
