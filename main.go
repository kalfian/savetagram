package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/kalfian/savetagram/downloader"
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
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Mohon tunggu sebentar...")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		linkMedia, typeMedia := scrapper.GetUrlInstagram(update.Message.Text)

		if linkMedia == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Gagal memperoleh data...")

			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}

		fmt.Println("-----------------------------------------")
		fmt.Println(linkMedia)
		fmt.Println("-----------------------------------------")

		fileName, err := downloader.DownloadFile(linkMedia, bot)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if typeMedia == scrapper.VIDEO {

			videoMsg := tgbotapi.NewVideoUpload(update.Message.Chat.ID, fileName)
			videoMsg.ReplyToMessageID = update.Message.MessageID
			_, err = bot.Send(videoMsg)
			if err != nil {
				log.Println(err.Error())
			}
			os.Remove(fileName)
		} else if typeMedia == scrapper.PHOTO {
			photoMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, fileName)
			photoMsg.ReplyToMessageID = update.Message.MessageID
			bot.Send(photoMsg)
			os.Remove(fileName)
		}
	}

	// fmt.Println(getUrlInstagram("https://www.instagram.com/p/CRyQhvqhQm7/?utm_source=ig_web_copy_link"))

}
