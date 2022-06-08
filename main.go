package main

import (
	"log"

	myBot "github.com/fengxxc/gakki_say/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// r := gin.Default()
	// r.Run(":1988")

	myBot.FetchTask(func(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		userUsername := update.Message.From.UserName
		userText := update.Message.Text
		log.Printf("[%s] %s", userUsername, userText)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if userText == "hello" {
			msg.Text = "hello, I'm Gakki~"
		} else {
			msg.Text = "(⊙o⊙)？"
		}
		bot.Send(msg)

	})
}
