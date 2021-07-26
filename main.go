package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/kalfian/savetagram/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

		linkMedia, typeMedia := getUrlInstagram(update.Message.Text)

		if linkMedia == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Gagal memperoleh data...")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}

		if typeMedia == VIDEO {
			videoMsg := tgbotapi.NewVideoUpload(update.Message.Chat.ID, linkMedia)
			videoMsg.ReplyToMessageID = update.Message.MessageID
			bot.Send(videoMsg)
		} else if typeMedia == PHOTO {
			photoMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, linkMedia)
			photoMsg.ReplyToMessageID = update.Message.MessageID
			bot.Send(photoMsg)
		}
	}

}

var (
	VIDEO = 1
	PHOTO = 2
)

func getUrlInstagram(url string) (string, int) {
	link := ""
	typeLink := 0

	var wg sync.WaitGroup
	wg.Add(1)

	c := colly.NewCollector(
		//colly.CacheDir("./_instagram_cache/"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		// fmt.Printf("%+v", e)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e, r.Request.URL, string(r.Body))
		link = ""
		wg.Done()
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Printf("%+v", string(r.Body))
		data := string(r.Body)
		delimiter := "window._sharedData = "
		dataSplited := strings.Split(data, delimiter)
		if len(dataSplited) > 0 {
			delimiter2 := ";</script>"
			splitAgain := strings.Split(dataSplited[1], delimiter2)

			// fmt.Printf("%+v", splitAgain[0])

			data := models.IgResponse{}

			err := json.Unmarshal([]byte(splitAgain[0]), &data)
			if err != nil {
				link = ""
				wg.Done()
			}

			log.Println(fmt.Sprintf("Data Mentah: %+v", splitAgain[0]))
			log.Println(fmt.Sprintf("Data Jadi: %+v", data))

			if len(data.EntryData.PostPage) > 0 {
				link = data.EntryData.PostPage[0].GraphQL.ShortcodeMedia.DisplayUrl
				typeLink = VIDEO

				if data.EntryData.PostPage[0].GraphQL.ShortcodeMedia.VideoUrl != "" {
					link = data.EntryData.PostPage[0].GraphQL.ShortcodeMedia.VideoUrl
					typeLink = PHOTO
				}
			}

		} else {
			log.Println("Data tidak sesuai")
		}

		wg.Done()
	})

	c.Visit(url)

	wg.Wait()

	return link, typeLink
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
