package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	case update.Message != nil &&
		update.Message.Entities != nil:
		resolveMessage(update.Message)

	case update.ChannelPost != nil:
		resolveChannelPost(update.ChannelPost)
	}
}

func resolveMessage(message *tgbotapi.Message) {

	switch {
	case regexp.MustCompile(`(?m)^\/list( \w+|$)`).MatchString(message.Text):
		adminId := os.Getenv(ADMIN_ID)
		if strconv.FormatInt(message.Chat.ID, 10) == adminId {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, getAllData()))
		}

	default:
		for _, entity := range *message.Entities {
			switch entity.Type {
			case HashtagType:
				resolveHashtagType(message, &entity)
			case CommandType:

			default:

			}
		}
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

func resolveHashtagType(message *tgbotapi.Message, entity *tgbotapi.MessageEntity) {
	hashtag := []rune(strings.ToLower(message.Text))[entity.Offset : entity.Length+entity.Offset]
	if _, exist := info[string(hashtag)]; exist {
		for _, channelId := range info[string(hashtag)] {
			msg := tgbotapi.NewMessage(channelId, message.Text)
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Ir a mensaje", fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)),
				),
			)

			_, err := bot.Send(msg)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func resolveCommandType(message *tgbotapi.Message, entity *tgbotapi.MessageEntity) {

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
		chat, err := bot.GetChat(tgbotapi.ChatConfig{
			ChatID: key,
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
		_, err := bot.SetWebhook(tgbotapi.NewWebhook(webhookPath))
		if err != nil {
			panic(err)
		}

		return bot.ListenForWebhook("/"), nil
	} else {
		return bot.GetUpdatesChan(tgbotapi.UpdateConfig{
			Offset:  0,
			Timeout: 0,
		})
	}
}
