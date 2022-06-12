package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	myTgBot "github.com/fengxxc/gakki_say/bot"
	policy "github.com/fengxxc/gakki_say/policy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// r := gin.Default()
	// r.Run(":1988")

	var config *Config = loadConfig("./config.json")

	var imgDef *policy.ImgDef = loadImgDef("./img/def.json")
	var symbolMaps policy.SymbolMaps = imgDef.GetMaps()
	// log.Printf("%+v", symbolMaps)

	myTgBot.FetchTask(config.TgBotToken, config.TgProxy, func(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		userUsername := update.Message.From.UserName
		userText := update.Message.Text
		log.Printf("[%s] %s", userUsername, userText)

		if update.Message.IsCommand() {
			// 处理命令
			myTgBot.CommmandHandler(bot, update.Message.Chat.ID, update.Message.Command())
		} else {
			// 处理信息
			myTgBot.UserTextHandler(bot, update.Message.Chat.ID, update.Message.MessageID, userText, &symbolMaps)
		}

	})

}

type Config struct {
	TgBotToken string `json:"tgBotToken"`
	TgProxy    string `json:"tgProxy"`
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

func loadImgDef(path string) *policy.ImgDef {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load img def file failed: ", err)
	}
	var imgDef *policy.ImgDef = &policy.ImgDef{}
	err = json.Unmarshal(f, &imgDef)
	if err != nil {
		log.Panicln("decode img def file failed: ", string(f), err)
	}
	// log.Printf("img def: %#v", imgDef)
	return imgDef
}
