package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	bot, err := tgbotapi.NewBotAPI("5465060326:AAGXybvWcExpT-RolQenge3PbVcJQx-mVm0")
	if err != nil {
		log.Printf("create telegram bot failed:%+v", err)
	}

	bot.Debug = true

	go func() {
		chatID := os.Getenv("CHAT_ID")
		chatIDValue, _ := strconv.ParseInt(chatID, 10, 64)
		msg := tgbotapi.NewMessage(chatIDValue, "你好呀")

		if _, err := bot.Send(msg); err != nil {
			log.Printf("send message to bot failed:%+v", err)
		}
	}()

	currentTime := time.Now().Format(time.RFC850)
	fmt.Fprintf(w, currentTime)
}
