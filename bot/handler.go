package bot

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	mapset "github.com/deckarep/golang-set"
	policy "github.com/fengxxc/gakki_say/policy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CallbackQueryHandler(bot *tgbotapi.BotAPI, chatId int64, messageId int, replyMessageId int, callbackQueryId string, callbackQueryData string, imgDir embed.FS, fontDir embed.FS) {
	// Respond to the callback query, telling Telegram to show the user
	// a message with the data received.
	/* callback := tgbotapi.NewCallback(callbackQueryId, callbackQueryData)
	log.Println(callback.Text)
	callback.Text = "You pick " + callback.Text
	if _, err := bot.Request(callback); err != nil {
		log.Println(err)
	} */

	// And finally, send a message containing the data received.
	switch callbackQueryData {
	case "some_case":
		// msg := tgbotapi.NewEditMessageTextAndMarkup(chatId, messageId, "几个栗子~ \n", inlineKeyboard)
		replyKB := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("🤬 谁放屁了？！"),
				tgbotapi.NewKeyboardButton("👀 让我康康~"),
				tgbotapi.NewKeyboardButton("☝️ 一定是你！"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("😊 萌混过关~"),
				tgbotapi.NewKeyboardButton("👐 卖萌禁止！"),
				tgbotapi.NewKeyboardButton("😞 心里委屈..."),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("😠 超级生气！"),
				tgbotapi.NewKeyboardButton("✌️ 宣告胜利~"),
				tgbotapi.NewKeyboardButton("[关闭]"),
			),
		)
		// msg := tgbotapi.NewEditMessageText(chatId, messageId, START_TEXT)
		msg := tgbotapi.NewMessage(chatId, "几个栗子~ 点击按钮快速发送。点击[关闭]退出。 \n")
		msg.ReplyMarkup = replyKB
		bot.Send(msg)
		return
	case "push_🏀":
		dice := tgbotapi.NewDiceWithEmoji(chatId, "🏀 (1-5)")
		// dice := tgbotapi.NewDice(chatId)
		dice.ReplyToMessageID = replyMessageId
		// log.Printf("push dice is: %+v\n", dice)
		sendMsg, err := bot.Send(dice)
		if err != nil {
			return
		}
		myVal := sendMsg.Dice.Value
		customerVal := sendMsg.ReplyToMessage.Dice.Value
		var rereText string
		if myVal > customerVal {
			rereText = fmt.Sprintf("你%d，我%d，我赢了~", customerVal, myVal)
		} else if myVal < customerVal {
			rereText = fmt.Sprintf("你%d，我%d，你赢了~", customerVal, myVal)
		} else {
			rereText = fmt.Sprintf("你%d，我%d，平局~", customerVal, myVal)
		}
		sendReply(bot, chatId, replyMessageId, policy.Reply{Type: policy.Text, Body: []byte(rereText)})
	}

}

var START_TEXT string = "初次见面，请多指教，我是图文并茂的gakki_say~ \n" +
	"你可以使用我生成带文字的Gakki图片 \n" +
	"我会根据emoji选择相应的Gakki图片并合成文字返回\n\n" +
	"具体方法是：\n" +
	"  发送 `emoji 你的文字`（注意emoji后有空格哦）\n" +
	"在群组里使用：\n" +
	"  1、加我进群，然后提拔我为管理员\n" +
	"  2、发送 `@gakki_say_bot emoji 你的文字`，没错，先at我，我才理你 \n" +
	"现在，输入 '👍 元气' 试试看~ \n\n" +
	"或者点击 '举个栗子' 快速体验\n" +
	"贡献代码或素材请点击 '项目地址'\n"

func CommmandHandler(bot *tgbotapi.BotAPI, chatId int64, command string, imgDir embed.FS, fontDir embed.FS) {
	var reply policy.Reply = policy.Reply{Type: policy.Failed, Body: []byte("")}
	switch command {
	case "start":
		msg := tgbotapi.NewMessage(chatId, START_TEXT)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("举个栗子", "some_case"),
				tgbotapi.NewInlineKeyboardButtonURL("项目地址", "https://github.com/fengxxc/gakki_say"),
			),
		)
		// msg.ParseMode = "Markdown"
		msg.Entities = []tgbotapi.MessageEntity{
			{Type: "code", Offset: policy.UnicodeIndex(START_TEXT, "emoji 你的文字"), Length: utf8.RuneCountInString("emoji 你的文字")},
			{Type: "code", Offset: policy.UnicodeIndex(START_TEXT, "@gakki_say_bot emoji 你的文字"), Length: utf8.RuneCountInString("@gakki_say_bot emoji 你的文字")},
			{Type: "code", Offset: policy.UnicodeIndex(START_TEXT, "👍 元气"), Length: utf8.RuneCountInString("👍 元气")},
		}
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
		img, err := policy.ImgWriteText(fileName, "pang?", policy.DrawStringConfig{
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

func DiceHandler(bot *tgbotapi.BotAPI, chatId int64, messageId int, dice *tgbotapi.Dice) {
	var msgText string
	var replyMap map[string][]string = make(map[string][]string)
	replyMap["🏀"] = []string{"没进，菜得抠脚", "兜兜转转，然而没进", "啊咧，非……非静止画面？", "篮网下面开口剪大点啊八嘎！", "好耶~算你投进了~"}
	replyMap["⚽"] = []string{"国足附体，再接再厉", "你就蹭蹭，不进去，嗯", "今天守门员请假，便宜你了~", "进了，角度刁钻~", "这也能进！好棒棒~"}
	if len(replyMap[dice.Emoji]) > 0 {
		msgText = replyMap[dice.Emoji][dice.Value-1] + " (" + strconv.Itoa(dice.Value) + ")"
	} else {
		msgText = strconv.Itoa(dice.Value)
	}
	var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("俺也来一个", "push_"+dice.Emoji),
		),
	)
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ReplyMarkup = inlineKeyboard
	msg.ReplyToMessageID = messageId
	bot.Send(msg)
	// sendReply(bot, chatId, messageId, policy.Reply{Type: policy.Text, Body: []byte(msgText)})
}

func UserTextHandler(bot *tgbotapi.BotAPI, chatId int64, chatType string, messageId int, replyMessageId int, userText string, symbolMaps *policy.SymbolMaps, imgDir embed.FS, fontDir embed.FS) {
	// selfBotName := "gakki_say_bot"
	// 在群、频道中，@我，我才会回应；私聊则不用
	/* if chatType != "private" {
		if !strings.HasPrefix(userText, "@"+selfBotName) {
			// 此处可偷听群聊~
			return
		}
		// @我了，remove @selfBotName
		userText = userText[len("@"+selfBotName):]
	} */

	switch userText {
	case "[关闭]":
		msg := tgbotapi.NewMessage(chatId, "关闭")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		sendMsg, _ := bot.Send(msg)
		bot.Send(tgbotapi.NewDeleteMessage(chatId, messageId))
		bot.Send(tgbotapi.NewDeleteMessage(chatId, sendMsg.MessageID))
		return
	}
	var reply policy.Reply = policy.UserText(userText, symbolMaps, imgDir, fontDir)
	msgId := messageId
	if replyMessageId != -1 {
		msgId = replyMessageId
	}
	sendReply(bot, chatId, msgId, reply)
}

func InlineQueryHandler(bot *tgbotapi.BotAPI, inlineQueryId string, query string, fromId int64, symbolMaps *policy.SymbolMaps) {
	var imgBaseURL string = "https://raw.githubusercontent.com/fengxxc/gakki_say/master/img/"
	var inlineQueryResults []interface{} = []interface{}{}
	query = strings.Trim(query, " ")
	emoji := strings.Split(query, " ")[0]
	if symbolMaps.ContainsEmoji(emoji) {
		var imgNameSet mapset.Set = symbolMaps.EmojiMap[emoji]
		imgNameSet.Each(func(item interface{}) bool {
			var imgName string = item.(string)
			var imgURL string = imgBaseURL + url.PathEscape(imgName)
			var resPhoto = tgbotapi.NewInlineQueryResultPhoto(md5Encode(imgName), imgURL)
			resPhoto.ThumbURL = imgURL
			inlineQueryResults = append(inlineQueryResults, resPhoto)
			return false
		})
	}
	if len(inlineQueryResults) == 0 {
		return
	}
	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryId,
		IsPersonal:    true,
		CacheTime:     5,
		Results:       inlineQueryResults,
	}
	bot.Request(inlineConfig)
	// bot.Send(inlineConfig)
}

func sendReply(bot *tgbotapi.BotAPI, chatId int64, messageId int, reply policy.Reply) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatId, "")
	if messageId != -1 {
		msg.ReplyToMessageID = messageId
	}
	var returnMsg tgbotapi.Message
	var err error
	if reply.Type == policy.Failed {
		msg.Text = "不要发令我困扰的东西哦……"
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

func md5Encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
