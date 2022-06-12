package bot

import (
	"log"

	policy "github.com/fengxxc/gakki_say/policy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommmandHandler(bot *tgbotapi.BotAPI, chatId int64, command string) {
	msg := tgbotapi.NewMessage(chatId, "")
	var reply policy.Reply = policy.Command(command)
	sendReply(bot, msg, chatId, reply)
}

func UserTextHandler(bot *tgbotapi.BotAPI, chatId int64, messageId int, userText string, symbolMaps *policy.SymbolMaps) {
	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("1"),
			tgbotapi.NewKeyboardButton("2"),
			tgbotapi.NewKeyboardButton("3"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("4"),
			tgbotapi.NewKeyboardButton("5"),
			tgbotapi.NewKeyboardButton("6"),
		),
	)

	var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4", "4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6", "6"),
		),
	)

	msg := tgbotapi.NewMessage(chatId, userText)
	msg.ReplyToMessageID = messageId

	switch userText {
	case "open":
		msg.ReplyMarkup = numericKeyboard
		bot.Send(msg)
		return
	case "close":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		bot.Send(msg)
		return
	case "open_inline":
		msg.ReplyMarkup = inlineKeyboard
		bot.Send(msg)
		return
	}
	var reply policy.Reply = policy.UserText(userText, symbolMaps)
	sendReply(bot, msg, chatId, reply)
}

func sendReply(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig, chatId int64, reply policy.Reply) {
	if reply.Type == policy.Failed {
		msg.Text = "吖白，大脑一片空白……"
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	} else if reply.Type == policy.Text {
		msg.Text = string(reply.Body)
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	} else if reply.Type == policy.Image {
		file := tgbotapi.FileBytes{
			Name:  "image.jpg",
			Bytes: reply.Body,
		}
		photo := tgbotapi.NewPhoto(chatId, file)
		if _, err := bot.Send(photo); err != nil {
			log.Println(err)
		}
	}
}
