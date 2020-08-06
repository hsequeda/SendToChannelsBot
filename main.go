package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	PORT         = "PORT"
	BOT_TOKEN    = "BOT_TOKEN"
	WEBHOOK_PATH = "WEBHOOK_PATH"
	HashtagType  = "hashtag"
	CommandType  = "command"
	ADMIN_ID     = "ADMIN_ID"
)

var bot *tgbotapi.BotAPI

var db *PostgresDatabase

var info map[string][]int64
var info2 map[string][]int64

func getCapasity() int {
	size := 0
	for hashtag, chIdList := range info {
		size += int(unsafe.Sizeof(hashtag))
		for _, chId := range chIdList {
			size += int(unsafe.Sizeof(chId))
		}
	}
	return size
}

func init() {
	var err error
	db, err = InitDb()
	if err != nil {
		panic(err)
	}

	info, err = db.List()
	if err != nil {
		panic(err)
	}
	info2 = info
	fmt.Printf("%#v bits\n", getCapasity())
	botToken := os.Getenv(BOT_TOKEN)
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	port := os.Getenv(PORT)
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func main() {
	fmt.Println("Started")
	updates, err := getUpdateCh()
	if err != nil {
		panic(err)
	}
	for update := range updates {
		resolveUpdate(&update)
	}
}

func resolveUpdate(update *tgbotapi.Update) {
	b, _ := json.Marshal(update)
	fmt.Printf("%#v\n\n", string(b))

	switch {
	case update.Message != nil:
		resolveMessage(update.Message)

	case update.ChannelPost != nil:
		resolveChannelPost(update.ChannelPost)
	}
}

func resolveMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		if err := handleCommand(message); err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	if message.Entities != nil || message.CaptionEntities != nil {
		resolveHashtagType(message)
	}
}

func resolveChannelPost(channelPost *tgbotapi.Message) {
	// if channelPost.Entities != nil {
	// 	for _, entity := range *channelPost.Entities {
	// 		if entity.Type == HashtagType {
	// 			hashtag := channelPost.Text[entity.Offset : entity.Offset+entity.Length]
	// 			if hashtag[0] == ' ' {
	// 				hashtag = channelPost.Text[entity.Offset+1 : entity.Offset+entity.Length+1]
	// 			}

	// 			if _, exist := info[hashtag]; !exist {
	// 				info[hashtag] = append(info[hashtag], channelPost.Chat.ID)
	// 				_, err := bot.Send(tgbotapi.NewMessage(channelPost.Chat.ID, fmt.Sprintf("Added hashtag: %s", hashtag)))
	// 				if err != nil {
	// 					fmt.Println(err.Error())
	// 				}
	// 			} else {
	// 				exist := false
	// 				for _, channelId := range info[hashtag] {
	// 					if channelPost.Chat.ID == channelId {
	// 						exist = true
	// 						break
	// 					}
	// 				}

	// 				if !exist {
	// 					info[hashtag] = append(
	// 						info[hashtag], channelPost.Chat.ID)
	// 					_, err := bot.Send(tgbotapi.NewMessage(channelPost.Chat.ID, fmt.Sprintf("Added hashtag: %s", hashtag)))
	// 					if err != nil {
	// 						fmt.Println(err.Error())
	// 					}
	// 				}
	// 			}
	// 		}

	// 		if err := db.Update(info); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	}
	// }
}

func resolveHashtagType(message *tgbotapi.Message) {
	if message.Text != "" {
		getMessageFromText(message)
	}

	if message.Caption != "" {
		getMessageFromCaption(message)
	}
}

func getAllData() string {
	result := map[int64][]string{}
	for key, value := range info {
		for _, value2 := range value {
			result[value2] = append(result[value2], key)
		}
	}

	var result2 string
	for key, value := range result {
		chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
			ChatConfig: tgbotapi.ChatConfig{ChatID: key},
		})

		if err != nil {
			fmt.Println(err.Error())
		}

		result2 += fmt.Sprintf("%s-> %s \n", chat.UserName, strings.Join(value, ", "))
	}

	return result2
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
	hashtagList := make([]string, 0)
	for _, entity := range entities {
		hashtagList = append(hashtagList, string([]rune(strings.ToLower(text))[entity.Offset:entity.Length+entity.Offset]))
	}

	return hashtagList
}

func getChannelList(hashtagList []string) []int64 {
	channelList := make([]int64, 0)
	channelFilter := make(map[int64]struct{})
	for _, hashtag := range hashtagList {
		// Get list of channels subscribed to a hashtag
		if subsChannels, ok := info[hashtag]; ok {
			for _, ch := range subsChannels {
				if _, ok := channelFilter[ch]; !ok {
					channelList = append(channelList, ch)
					channelFilter[ch] = struct{}{}
				}
			}
		}
	}

	return channelList
}

func handleCommand(message *tgbotapi.Message) error {
	switch {
	case regexp.MustCompile(`(?m)^\/list( \w+|$)`).MatchString(message.Text):
		adminId := os.Getenv(ADMIN_ID)
		if strconv.FormatInt(message.Chat.ID, 10) == adminId {
			_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, getAllData()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getMessageFromText(message *tgbotapi.Message) {
	hashtagList := getHashtagList(message.Text, message.Entities)
	channelIdList := getChannelList(hashtagList)
	for _, chId := range channelIdList {
		toSend := tgbotapi.NewMessage(chId, fmt.Sprintf("%s\n%s", message.Text, getRefLink(message.From)))
		toSend.ParseMode = "html"
		toSend.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
			),
		)

		msg, err := bot.Send(toSend)
		if err != nil {
			fmt.Println(err.Error())
		}
		b, _ := json.Marshal(msg)
		fmt.Println(string(b))
	}
}

func getMessageFromCaption(message *tgbotapi.Message) {
	hashtagList := getHashtagList(message.Caption, message.CaptionEntities)
	channelIdList := getChannelList(hashtagList)
	refLink := getRefLink(message.From)
	for _, chId := range channelIdList {
		var toSend tgbotapi.Chattable
		if message.Photo != nil {
			photoConfig := tgbotapi.NewPhotoShare(chId, message.Photo[0].FileID)
			photoConfig.Caption = fmt.Sprintf("%s\n%s", message.Caption, refLink)
			photoConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
				),
			)
			toSend = photoConfig
		}

		if message.Audio != nil {
			audioConfig := tgbotapi.NewAudioShare(chId, message.Audio.FileID)
			audioConfig.Caption = fmt.Sprintf("%s\n%s", message.Caption, refLink)
			audioConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
				),
			)
			toSend = audioConfig
		}

		if message.Video != nil {
			videoConfig := tgbotapi.NewVideoShare(chId, message.Video.FileID)
			videoConfig.Caption = fmt.Sprintf("%s\n%s", message.Caption, refLink)
			videoConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
				),
			)
			toSend = videoConfig
		}

		if message.Document != nil {
			documentConfig := tgbotapi.NewDocumentShare(chId, message.Document.FileID)
			documentConfig.Caption = fmt.Sprintf("%s\n%s", message.Caption, refLink)
			documentConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
				),
			)
			toSend = documentConfig
		}

		_, err := bot.Send(toSend)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func getRefLink(user *tgbotapi.User) string {
	var name string
	name = user.UserName
	if user.UserName == "" {
		name = user.FirstName
	}
	return fmt.Sprintf("\n <a href=\"tg://user?id=%d\">Hablar con autor👤(%s)</a> ", user.ID, name)
}
