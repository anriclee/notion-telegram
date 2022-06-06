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

func Handler(w http.ResponseWriter, r *http.Request) {
	body, _ := r.GetBody()
	bytes, _ := io.ReadAll(body)
	msg := string(bytes)
	log.Printf("receive msg from telegram,request body is:%+v", msg)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Printf("create telegram bot failed:%+v", err)
	}

	bot.Debug = true

	go func() {
		chatID := os.Getenv("CHAT_ID")
		chatIDValue, _ := strconv.ParseInt(chatID, 10, 64)
		msg := tgbotapi.NewMessage(chatIDValue, msg)

		if _, err := bot.Send(msg); err != nil {
			log.Printf("send message to bot failed:%+v", err)
		}
	}()

	currentTime := time.Now().Format(time.RFC850)
	fmt.Fprintf(w, currentTime)
}
