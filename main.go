package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/kalfian/savetagram/scrapper"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	// go http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	for update := range updates {
		scrapper.Handle(update, bot)
	}

	// fmt.Println(getUrlInstagram("https://www.instagram.com/p/CRyQhvqhQm7/?utm_source=ig_web_copy_link"))

}
