package scrapper

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	"github.com/kalfian/savetagram/downloader"
	"github.com/kalfian/savetagram/models"
)

var (
	VIDEO = 1
	PHOTO = 2
)

func GetUrlInstagram(url string) (string, int) {
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
		// r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Referer", "https://www.instagram.com/"+"2626525925817652203")
		if r.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", r.Ctx.Get("gis"), r.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			r.Headers.Set("X-Instagram-GIS", gisHash)
		}
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

			if len(data.EntryData.PostPage) > 0 {
				link = data.EntryData.PostPage[0].GraphQL.ShortcodeMedia.DisplayUrl
				typeLink = PHOTO

				if data.EntryData.PostPage[0].GraphQL.ShortcodeMedia.VideoUrl != "" {
					link = data.EntryData.PostPage[0].GraphQL.ShortcodeMedia.VideoUrl
					typeLink = VIDEO
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

func Handle(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	if update.Message == nil { // ignore any non-Message Updates
		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Mohon tunggu sebentar...")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	linkMedia, typeMedia := GetUrlInstagram(update.Message.Text)

	if linkMedia == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Gagal memperoleh data...")

		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
		return nil
	}

	fmt.Println("-----------------------------------------")
	fmt.Println(linkMedia)
	fmt.Println("-----------------------------------------")

	fileName, err := downloader.DownloadFile(linkMedia, bot)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if typeMedia == VIDEO {

		videoMsg := tgbotapi.NewVideoUpload(update.Message.Chat.ID, fileName)
		videoMsg.ReplyToMessageID = update.Message.MessageID
		_, err = bot.Send(videoMsg)
		if err != nil {
			log.Println(err.Error())
		}
		os.Remove(fileName)
	} else if typeMedia == PHOTO {
		photoMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, fileName)
		photoMsg.ReplyToMessageID = update.Message.MessageID
		bot.Send(photoMsg)
		os.Remove(fileName)
	}

	return nil
}
