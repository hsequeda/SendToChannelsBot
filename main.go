package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"os"
)

var bot *tgbotapi.BotAPI

func init() {
	botToken := os.Getenv("BOT_TOKEN")
	var err error
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
		if update.Message.Entities != nil {
			for _, entity := range *update.Message.Entities {
				if entity.Type == "hashtag" {
					_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text[entity.Offset:entity.Length+entity.Offset]))
					if err != nil {
						println(err.Error())
					}
				}
			}
		}
	}
}
