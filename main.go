package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	myTgBot "github.com/fengxxc/gakki_say/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// r := gin.Default()
	// r.Run(":1988")

	config := loadConfig("./config.json")

	myTgBot.FetchTask(config.TgBotToken, func(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		userUsername := update.Message.From.UserName
		userText := update.Message.Text
		log.Printf("[%s] %s", userUsername, userText)

		if update.Message.IsCommand() {
			// 处理命令
			myTgBot.CommmandHandler(bot, update.Message.Chat.ID, update.Message.Command())
		} else {
			// 处理信息
			myTgBot.UserTextHandler(bot, update.Message.Chat.ID, update.Message.MessageID, userText)
		}

	})
}

type Config struct {
	TgBotToken string `json:"tgBotToken"`
}

func loadConfig(path string) *Config {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config file failed: ", err)
	}
	var config *Config = &Config{}
	err = json.Unmarshal(f, config)
	if err != nil {
		log.Panicln("decode config file failed:", string(f), err)
	}
	return config
}
