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

var client = http.Client{Timeout: 5 * time.Minute}

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

		content, err := reqQRCode()
		if err != nil {
			log.Printf("request qr code failed:%+v", err)
			return
		}

		png, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Printf("encode qr failed:%+v", err)
			return
		}

		msg := tgbotapi.NewPhoto(chatIDValue, tgbotapi.FileBytes{"code", png})
		if _, err := bot.Send(msg); err != nil {
			log.Printf("send message to bot failed:%+v", err)
		}
	}()

	fmt.Fprintf(w, "receive cmd ok")
}

func reqQRCode() (string, error) {
	req, err := http.NewRequest("POST", "http://doorcloud.sohochina.com/rest/sohoweCharTect/getOwnerQrCode", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	req.Form.Add("userLinglingid", "00EEF073")
	req.Form.Add("supportControl", "0")
	req.Form.Add("jurId", "278")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("code not ok:%v", response.StatusCode)
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
