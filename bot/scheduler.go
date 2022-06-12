package bot

import (
	"log"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FetchTask(botToken string, proxy string, updateCallback func(*tgbotapi.BotAPI, tgbotapi.Update)) {
	var myClient *http.Client = &http.Client{}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			log.Panic(err)
		}
		myClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}
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
