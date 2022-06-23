package bot

import (
	"embed"
	"log"

	policy "github.com/fengxxc/gakki_say/policy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CallbackQueryHandler(bot *tgbotapi.BotAPI, chatId int64, callbackQueryId string, callbackQueryData string, imgDir embed.FS, fontDir embed.FS) {
	// Respond to the callback query, telling Telegram to show the user
	// a message with the data received.
	/* callback := tgbotapi.NewCallback(callbackQueryId, callbackQueryData)
	log.Println(callback.Text)
	callback.Text = "You pick " + callback.Text
	if _, err := bot.Request(callback); err != nil {
		log.Println(err)
	} */

	// And finally, send a message containing the data received.
	if callbackQueryData == "random_case" {
		randomText := "🤬 八嘎，谁放屁了？！"
		callback := tgbotapi.NewCallback(callbackQueryId, "you send '"+randomText+"'")
		if _, err := bot.Request(callback); err != nil {
			log.Println(err)
		}
		userMsg := sendReply(bot, chatId, -1, policy.Reply{
			Type: policy.Text,
			Body: []byte(randomText),
		})

		fileName := "img/逃げ恥_09.003633.146.png"
		img, err := policy.ImgWriteTextDefault(fileName, randomText, imgDir, fontDir)
		if err != nil {
			log.Println(err)
			return
		}
		var reply policy.Reply = policy.Reply{
			Type: policy.Image,
			Body: policy.ImgToBytes(img, policy.GetImgTypeByFileName(fileName)),
		}
		sendReply(bot, chatId, userMsg.MessageID, reply)
	}

}

func CommmandHandler(bot *tgbotapi.BotAPI, chatId int64, command string, imgDir embed.FS, fontDir embed.FS) {
	var reply policy.Reply = policy.Reply{Type: policy.Failed, Body: []byte("")}
	switch command {
	case "start":
		// reply.Type = policy.Text
		// reply.Body = []byte("初次见面，请多指教，我是图文并茂的Gakki~")
		msg := tgbotapi.NewMessage(chatId, "初次见面，请多指教，我是图文并茂的gakki_say~ \n"+
			"你可以使用我生成带文字的Gakki图片。 \n\n"+
			"具体方法是：聊天框中输入 'emoji表情[空格]展现的文字'\n"+
			"我会根据emoji选择相应的Gakki图片并合成文字返回\n"+
			"贡献代码或素材请点击“项目地址”\n\n"+
			"现在，输入'👍 元气'试试看~ ",
		)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("项目地址", "https://github.com/fengxxc/gakki_say"),
				tgbotapi.NewInlineKeyboardButtonData("随机示例", "random_case"),
			),
		)
		bot.Send(msg)
		return
	case "help":
		reply.Type = policy.Text
		reply.Body = []byte("我还没想好怎么帮你")
	case "settings":
		reply.Type = policy.Text
		reply.Body = []byte("这个功能还没做好……再等等")
	case "ping":
		fileName := "img/pingpang.jpg"
		img, err := policy.ImgWriteText(fileName, "pang~", policy.DrawStringConfig{
			Ax:          0.5,
			Ay:          0.5,
			FontFamily:  "SIMYOU.TTF",
			TextBgColor: &policy.RGBA{R: 255, G: 204, B: 255, A: 89},
		}, imgDir, fontDir)
		if err != nil {
			log.Println(err)
			return
		}
		reply.Type = policy.Image
		reply.Body = policy.ImgToBytes(img, policy.GetImgTypeByFileName(fileName))
	}

	// var reply policy.Reply = policy.Command(command)
	sendReply(bot, chatId, -1, reply)
}

func UserTextHandler(bot *tgbotapi.BotAPI, chatId int64, messageId int, userText string, symbolMaps *policy.SymbolMaps, imgDir embed.FS, fontDir embed.FS) {
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

	switch userText {
	case "open":
		msg := tgbotapi.NewMessage(chatId, "open~")
		msg.ReplyMarkup = numericKeyboard
		bot.Send(msg)
		return
	case "close":
		msg := tgbotapi.NewMessage(chatId, "close~")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		bot.Send(msg)
		return
	case "open_inline":
		msg := tgbotapi.NewMessage(chatId, "open_inline~")
		msg.ReplyMarkup = inlineKeyboard
		bot.Send(msg)
		return
	}
	var reply policy.Reply = policy.UserText(userText, symbolMaps, imgDir, fontDir)
	sendReply(bot, chatId, messageId, reply)
}

func sendReply(bot *tgbotapi.BotAPI, chatId int64, messageId int, reply policy.Reply) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatId, "")
	if messageId != -1 {
		msg.ReplyToMessageID = messageId
	}
	var returnMsg tgbotapi.Message
	var err error
	if reply.Type == policy.Failed {
		msg.Text = "吖白，大脑一片空白……"
		returnMsg, err = bot.Send(msg)
	} else if reply.Type == policy.Text {
		msg.Text = string(reply.Body)
		returnMsg, err = bot.Send(msg)
	} else if reply.Type == policy.Image {
		file := tgbotapi.FileBytes{
			Name:  "image.jpg",
			Bytes: reply.Body,
		}
		photo := tgbotapi.NewPhoto(chatId, file)
		if messageId != -1 {
			photo.ReplyToMessageID = messageId
		}
		returnMsg, err = bot.Send(photo)
	}
	if err != nil {
		log.Println(err)
	}
	return returnMsg
}
