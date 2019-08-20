package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error

	bot, err = linebot.New(
			os.Getenv("CHANNEL_SECRET"),
			os.Getenv("CHANNEL_TOKEN"),
		)

	log.Print("Bot:", bot, "err:", err)

	http.HandleFunc("/callback", callbackHandler)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)

	log.Print("server start up:", addr)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				// free account
				if quota.Value >= 500 {
					log.Println("Quota not enough!")
				}
				err = msgFunc(event.ReplyToken, message.Text)
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func msgFunc(token, msg string) error {
	_, err := bot.ReplyMessage(token, linebot.NewTextMessage(msg)).Do()
	return err
}
