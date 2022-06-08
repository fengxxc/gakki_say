package bot

import (
	"log"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FetchTask(updateCallback func(*tgbotapi.BotAPI, tgbotapi.Update)) {
	botToken := "5484491783:AAGeic7Cex9bcjXx2dXfi6r0nZlC_gYs-x4"

	// bot, err := tgbotapi.NewBotAPI(botToken)
	proxyUrl, err := url.Parse("socks5://127.0.0.1:1070") //设置代理http或sock5
	if err != nil {
		log.Panic(err)
	}
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}} //使用代理
	bot, err := tgbotapi.NewBotAPIWithClient(botToken, tgbotapi.APIEndpoint, myClient)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			updateCallback(bot, update)
		}
	}
}
