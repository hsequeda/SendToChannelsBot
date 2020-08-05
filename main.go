package main

import (
	"fmt"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	PORT         = "PORT"
	BOT_TOKEN    = "BOT_TOKEN"
	WEBHOOK_PATH = "WEBHOOK_PATH"
	HashtagType  = "hashtag"
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

	webHookPath := os.Getenv(WEBHOOK_PATH)
	resp, err := bot.SetWebhook(tgbotapi.NewWebhook(webHookPath))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", resp)

	port := os.Getenv(PORT)
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func main() {
	fmt.Println("Started")
	for update := range bot.ListenForWebhook("/") {
		resolveUpdate(&update)
	}
}

func resolveUpdate(update *tgbotapi.Update) {
	switch {
	case update.Message != nil && update.Message.Entities != nil:
		resolveMessage(update.Message)
	case update.ChannelPost != nil:
		resolveChannelPost(update.ChannelPost)
	}
}

func resolveMessage(message *tgbotapi.Message) {
	for _, entity := range *message.Entities {
		if entity.Type == HashtagType {
			hashtag := message.Text[entity.Offset : entity.Length+entity.Offset]
			if hashtag[0] == ' ' {
				hashtag = message.Text[entity.Offset+1 : entity.Offset+entity.Length+1]
			}
			if _, exist := info[hashtag]; exist {
				for _, channelId := range info[hashtag] {
					_, err := bot.Send(tgbotapi.NewMessage(channelId, message.Text))
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			} else {
				if _, err := bot.Send(
					tgbotapi.NewMessage(
						message.Chat.ID,
						fmt.Sprintf("Hashtag: %s isn't added", hashtag),
					),
				); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

func resolveChannelPost(channelPost *tgbotapi.Message) {
	if channelPost.Entities != nil {
		for _, entity := range *channelPost.Entities {
			if entity.Type == HashtagType {
				hashtag := channelPost.Text[entity.Offset : entity.Offset+entity.Length]
				if hashtag[0] == ' ' {
					hashtag = channelPost.Text[entity.Offset+1 : entity.Offset+entity.Length+1]
				}

				if _, exist := info[hashtag]; !exist {
					info[hashtag] = append(info[hashtag], channelPost.Chat.ID)
					_, err := bot.Send(tgbotapi.NewMessage(channelPost.Chat.ID, fmt.Sprintf("Added hashtag: %s", hashtag)))
					if err != nil {
						fmt.Println(err.Error())
					}
				} else {
					exist := false
					for _, channelId := range info[hashtag] {
						if channelPost.Chat.ID == channelId {
							exist = true
							break
						}
					}

					if !exist {
						info[hashtag] = append(
							info[hashtag], channelPost.Chat.ID)
						_, err := bot.Send(tgbotapi.NewMessage(channelPost.Chat.ID, fmt.Sprintf("Added hashtag: %s", hashtag)))
						if err != nil {
							fmt.Println(err.Error())
						}
					}
				}
			}

			if err := db.Update(info); err != nil {
				fmt.Println(err)
			}
		}
	}
}
