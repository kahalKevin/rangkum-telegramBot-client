package main

import (
	"bytes"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("318149659:AAHaM5VuYYoMgjQM7rCDD8L42JbxK254b6o")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	rangkumAPI := "http://rangkum.herokuapp.com/rangkum318149659:AAHaM5VuYYoMgjQM7rCDD8L42JbxK254b6o"
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		reply := ""

		if strings.HasPrefix(update.Message.Text, "/about") {
			reply = "This is a bot to help you summarize article from any url, it's not really smart yet, only can handle english"
		} else if strings.HasPrefix(update.Message.Text, "/rangkum ") {
			reply = update.Message.Text[9:]
			var jsonStr = []byte(`{"url":"`+reply+`"}`)
			
			req, err := http.NewRequest("POST", rangkumAPI, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				log.Printf("[%s] %s", resp, err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			reply = string(body)
			if strings.HasPrefix(reply, "<!DOCTYPE HTML PUBLIC "){
				reply = "your url is invalid :)"
			}
		} else {
			reply = "You can type \n/rangkum [url]\n/about"
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}