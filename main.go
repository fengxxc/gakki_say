package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/fengxxc/gakki_say/api"
	myTgBot "github.com/fengxxc/gakki_say/bot"
	policy "github.com/fengxxc/gakki_say/policy"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/imdario/mergo"
)

//go:embed img
var imgDir embed.FS

//go:embed font
var fontDir embed.FS

func main() {

	var config *Config = loadConfig()
	// confMap := loadConfigToMap()
	args := os.Args
	argsMap := osArgs2Map(args, "--")
	log.Printf("%+v", argsMap)
	// mergo.MergeWithOverwrite(config, argsMap)
	mergo.Map(config, argsMap, mergo.WithOverwriteWithEmptyValue)
	log.Printf("%+v", config)

	var imgDef *policy.ImgDef = loadImgDef()
	var symbolMaps policy.SymbolMaps = imgDef.GetMaps()
	// log.Printf("%+v", symbolMaps)

	// make thumbnail
	makeThumbnail()

	r := gin.Default()
	api.Router(r, imgDir, fontDir)
	go func() {
		r.Run(config.ServerListen)
	}()

	myTgBot.FetchTask(config.TgBotToken, config.TgProxy, func(tgUpdType myTgBot.TgUpdType, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		switch tgUpdType {
		case myTgBot.Message:
			userUsername := update.Message.From.UserName
			userText := update.Message.Text
			log.Printf("[%s] %s", userUsername, userText)
			if update.Message.IsCommand() {
				// 处理命令
				myTgBot.CommmandHandler(bot, update.Message.Chat.ID, update.Message.Command(), imgDir, fontDir)
			} else if update.Message.Dice != nil {
				// 处理 骰子、飞镖、保龄球、篮球、足球、老虎机
				var dice *tgbotapi.Dice = update.Message.Dice
				myTgBot.DiceHandler(bot, update.Message.Chat.ID, update.Message.MessageID, dice)
			} else {
				replyMessageId := -1
				if replyToMessage := update.Message.ReplyToMessage; replyToMessage != nil {
					replyMessageId = replyToMessage.MessageID
				}
				// 处理文本
				myTgBot.UserTextHandler(bot, update.Message.Chat.ID, update.Message.Chat.Type, update.Message.MessageID, replyMessageId, userText, &symbolMaps, imgDir, fontDir)
			}

		case myTgBot.CallbackQuery:
			replyMessageId := -1
			if replyToMessage := update.CallbackQuery.Message.ReplyToMessage; replyToMessage != nil {
				replyMessageId = replyToMessage.MessageID
			}
			myTgBot.CallbackQueryHandler(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, replyMessageId, update.CallbackQuery.ID, update.CallbackQuery.Data, imgDir, fontDir)
		case myTgBot.InlineQuery:
			myTgBot.InlineQueryHandler(bot, update.InlineQuery.ID, update.InlineQuery.Query, update.InlineQuery.From.ID, &symbolMaps, config.Host, imgDir)
		}

	})

}

type Config struct {
	TgBotToken   string `json:"tgBotToken"`
	TgProxy      string `json:"tgProxy"`
	Host         string `json:"host"`
	ServerListen string `json:"serverListen"`
}

//go:embed config.json
var configFile embed.FS

func loadConfig() *Config {
	// configFile, err := ioutil.ReadFile(path)
	data, err := configFile.ReadFile("config.json")
	if err != nil {
		log.Panicln("load config file failed: ", err)
	}
	var config *Config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Panicln("decode config file failed:", string(data), err)
	}
	return config
}

//go:embed img/def.json
var defFile embed.FS

func loadImgDef() *policy.ImgDef {
	// f, err := ioutil.ReadFile(path)
	data, err := defFile.ReadFile("img/def.json")
	if err != nil {
		log.Panicln("load img def file failed: ", err)
	}
	var imgDef *policy.ImgDef = &policy.ImgDef{}
	err = json.Unmarshal(data, &imgDef)
	if err != nil {
		log.Panicln("decode img def file failed: ", string(data), err)
	}
	// log.Printf("img def: %#v", imgDef)
	return imgDef
}

func osArgs2Map(args []string, argPrefix string) map[string]interface{} {
	if args == nil || len(args) <= 1 {
		return nil
	}
	argMap := make(map[string]interface{})
	tempKey := ""
	for i := 1; i < len(args); i++ {
		if i%2 != 0 {
			tempKey = args[i][len(argPrefix):]
			continue
		}
		argMap[tempKey] = args[i]
	}
	return argMap
}

func makeThumbnail() {
	thumbDir, _ := os.Getwd()
	thumbDir = thumbDir + "/.cache/thumb"
	err := os.RemoveAll(thumbDir)
	if err != nil {
		log.Panicf("remove thumb dir failed: %s", err)
	}
	err = os.MkdirAll(thumbDir, 0755)
	if err != nil {
		log.Panicf("make thumb dir failed: %s", err)
	}
	fs.WalkDir(imgDir, "img", func(path string, d fs.DirEntry, err error) error {
		// log.Printf("walkDir: %s", path)
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() == "def.json" {
			return nil
		}
		log.Printf("dirEntry.Name: %+v", d.Name())
		file, err := imgDir.Open(path)
		if err != nil {
			log.Println("open file failed: ", err)
			return err
		}
		img, err := imaging.Decode(file)
		if err != nil {
			log.Println("decode file failed: ", err)
		}
		img = imaging.Resize(img, 0, 89, imaging.Lanczos)
		err = imaging.Save(img, thumbDir+"/"+d.Name(), imaging.JPEGQuality(80))
		if err != nil {
			log.Printf("save thumb failed: %s", err)
			return nil
		}
		return nil
	})
}
