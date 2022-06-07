package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBody struct {
	UpdatedID string  `json:"update_id"`
	Message   Message `json:"message"`
}

type Message struct {
	MessageID string `json:"message_id"`
	From      struct {
		ID    int  `json:"id"`
		IsBot bool `json:"is_bot"`
	} `json:"from"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("parse form values failed:%+v", err)
	}
	body, _ := r.GetBody()
	bytes, _ := io.ReadAll(body)
	msg := string(bytes)
	log.Printf("receive msg from telegram,request body is:%+v,id:%v", msg, r.FormValue("id"))

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Printf("create telegram bot failed:%+v", err)
	}

	bot.Debug = true

	go func() {
		chatID := os.Getenv("CHAT_ID")
		chatIDValue, _ := strconv.ParseInt(chatID, 10, 64)
		msg := tgbotapi.NewMessage(chatIDValue, msg)

		msg.text += r.FormValue("id")

		if _, err := bot.Send(msg); err != nil {
			log.Printf("send message to bot failed:%+v", err)
		}
	}()

	currentTime := time.Now().Format(time.RFC850)
	fmt.Fprintf(w, currentTime)
}
