package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI
var db = make(map[string][]string)

func init() {
	file, err := os.OpenFile("./data.json", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile("./data.json")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	if err = json.NewDecoder(file).Decode(&db); err != nil {
		fmt.Println(err.Error())
		db = make(map[string][]string, 0)
	}

	botToken := os.Getenv("BOT_TOKEN")
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	webHookPath := os.Getenv("WEBHOOK_PATH")
	println(webHookPath)
	resp, err := bot.SetWebhook(tgbotapi.NewWebhook(webHookPath))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", resp)

	port := os.Getenv("PORT")
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
		if entity.Type == "hashtag" {
			if _, exist := db[message.Text[entity.Offset:entity.Length+entity.Offset]]; exist {
				for _, channelId := range db[message.Text[entity.Offset:entity.Length+entity.Offset]] {
					id, err := strconv.ParseInt(channelId, 10, 64)
					if err != nil {
						fmt.Println(err.Error())
					}

					_, err = bot.Send(tgbotapi.NewMessage(id, message.Text))
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			} else {
				_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Hashtag: %s isn't added",
					message.Text[entity.Offset:entity.Length+entity.Offset])))
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

func resolveChannelPost(channelPost *tgbotapi.Message) {
	// adminId := os.Getenv("ADMIN_ID")
	// if string(channelPost.From.ID) == adminId &&
	if channelPost.Entities != nil {
		for _, entity := range *channelPost.Entities {
			if entity.Type == "hashtag" {
				if _, exist := db[channelPost.Text[entity.Offset:entity.Offset+entity.Length]]; !exist {
					db[channelPost.Text[entity.Offset:entity.Offset+entity.Length]] = make([]string, 0)
					db[channelPost.Text[entity.Offset:entity.Offset+entity.Length]] = append(db[channelPost.Text[entity.Offset:entity.Offset+entity.Length]],
						strconv.FormatInt(channelPost.Chat.ID, 10))
					_, err := bot.Send(tgbotapi.NewMessage(channelPost.Chat.ID, fmt.Sprintf("Added hashtag: %s", channelPost.Text[entity.Offset:entity.Offset+entity.Length])))
					if err != nil {
						fmt.Println(err.Error())
					}
				} else {
					exist := false
					for _, channelId := range db[channelPost.Text[entity.Offset:entity.Offset+entity.Length]] {
						if strconv.FormatInt(channelPost.Chat.ID, 10) == channelId {
							exist = true
							break
						}
					}

					if !exist {
						db[channelPost.Text[entity.Offset:entity.Offset+entity.Length]] = append(
							db[channelPost.Text[entity.Offset:entity.Offset+entity.Length]], strconv.FormatInt(channelPost.Chat.ID, 10))
						_, err := bot.Send(tgbotapi.NewMessage(channelPost.Chat.ID, fmt.Sprintf("Added hashtag: %s", channelPost.Text[entity.Offset:entity.Offset+entity.Length])))
						if err != nil {
							fmt.Println(err.Error())
						}
					}
				}
			}

			file, err := os.OpenFile("./data.json", os.O_CREATE|os.O_RDWR, 0777)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = json.NewEncoder(file).Encode(db)
			if err != nil {
				fmt.Println(err.Error())
			}

			err = json.NewDecoder(file).Decode(&db)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(db)

			file.Close()
		}
	}
}
