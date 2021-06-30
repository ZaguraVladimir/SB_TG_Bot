package main

import (
	"SB_TG_Bot/model"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	db := make(model.DB)

	bot, err := tgbotapi.NewBotAPI("КАКОЙ-ТО ТОКЕН")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		chatID := model.ChatID(update.Message.Chat.ID)

		db.AddID(chatID)

		message, err := model.NewMessage(update.Message.Text)
		if err != nil {
			sendMessage(bot, chatID, err.Error())
			continue
		}

		wallet := db[chatID]
		result := wallet.Processing(message)
		sendMessage(bot, chatID, result)
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID model.ChatID, textMessage string) {
	tgms := tgbotapi.NewMessage(chatID.ID(), textMessage)
	if _, err := bot.Send(tgms); err != nil {
		log.Println(err)
	}
}
