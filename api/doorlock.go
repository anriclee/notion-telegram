package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
)

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

		var png []byte
		png, err := qrcode.Encode("1cFW4h6xgXolhx7ewW447xyMUa0", qrcode.Medium, 256)
		if err != nil {
			log.Printf("encode qr failed:%+v", err)
		}

		msg := tgbotapi.NewPhoto(chatIDValue, tgbotapi.FileBytes{"code", png})

		if _, err := bot.Send(msg); err != nil {
			log.Printf("send message to bot failed:%+v", err)
		}
	}()

	currentTime := time.Now().Format(time.RFC850)
	fmt.Fprintf(w, currentTime)
}
