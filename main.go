package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mmcdole/gofeed"
)

var (
	tgToken = os.Getenv("TG_TOKEN_NEWS")

	// Define RSS feed URLs as constants
	autoFeedURL          = "https://3dnews.ru/auto/rss/"
	gadgetsFeedURL       = "https://3dnews.ru/gadgets/rss/"
	breakingFeedURL      = "https://3dnews.ru/breaking/rss"
	memeFeedURL          = "https://pikabu.ru/xmlfeeds.php?cmd=popular"
	tvboxFeedURL         = "https://www.reddit.com/r/AndroidTVBoxes/.rss"
	redditMemeFeedURL    = "https://www.reddit.com/r/pikabu/.rss"
	redditMobilePhotoUrl = "https://www.reddit.com/r/mobilephotography/.rss"
	anekdot              = "https://www.anekdot.ru/rss/export_bestday.xml"
	history              = "https://www.anekdot.ru/rss/export_o.xml"
	photoNews            = "http://lenta.ru/rss/photo"
	bashOrg              = "https://bashorg.org/rss.xml"
	calend               = "https://www.calend.ru/img/export/today-holidays.rss"
)

func main() {
	fmt.Println("Bot is starting...")

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Fatal(err)
	}

	// Configure update polling
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start", "/help":
			sendHelpMessage(bot, update.Message.Chat.ID)
		case "/today":
			sendMessage(bot, update.Message.Chat.ID, getNews(calend, 3))
		case "/photoNews":
			sendMessage(bot, update.Message.Chat.ID, randomMeme(photoNews))
		case "/history":
			sendMessage(bot, update.Message.Chat.ID, randomAnekdot(history))
		case "/anekdot":
			sendMessage(bot, update.Message.Chat.ID, randomAnekdot(anekdot))
		case "/rmp":
			sendMessage(bot, update.Message.Chat.ID, randomMeme(redditMobilePhotoUrl))
		case "/rmeme":
			sendMessage(bot, update.Message.Chat.ID, randomMeme(redditMemeFeedURL))
		case "/meme":
			sendMessage(bot, update.Message.Chat.ID, randomMeme(memeFeedURL))
		case "/tvbox":
			sendMessage(bot, update.Message.Chat.ID, getNews(tvboxFeedURL, 5))
		case "/breaking":
			sendMessage(bot, update.Message.Chat.ID, getNews(breakingFeedURL, 5))
		case "/auto":
			sendMessage(bot, update.Message.Chat.ID, getNews(autoFeedURL, 5))
		case "/gadgets":
			sendMessage(bot, update.Message.Chat.ID, getNews(gadgetsFeedURL, 5))
		case "/bash":
			sendMessage(bot, update.Message.Chat.ID, randomAnekdot(bashOrg))
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.DisableNotification = true
	//msg.DisableWebPagePreview = true

	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}

func randomAnekdot(url string) string {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Println("Error parsing meme feed:", err)
		return "Failed to fetch meme"
	}

	//rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(feed.Items))
	news := feed.Items[index]
	news.Description = strings.ReplaceAll(news.Description, "<br>", "\n")
	news.Description = strings.ReplaceAll(news.Description, "<br/>", "\n")
	news.Description = strings.ReplaceAll(news.Description, "<br />", "\n")
	return fmt.Sprintf("%s : %s", news.Title, news.Description)
}

func randomMeme(url string) string {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Println("Error parsing meme feed:", err)
		return "Failed to fetch meme"
	}

	//rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(feed.Items))
	news := feed.Items[index]
	return fmt.Sprintf("%s : %s", news.Title, news.Link)
}

func getNews(url string, numItems int) string {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Println("Error parsing feed:", err)
		return "Failed to fetch news"
	}

	var newsItems []string
	for i := 0; i < numItems && i < len(feed.Items); i++ {
		item := feed.Items[i]
		newsItems = append(newsItems, fmt.Sprintf("%s : %s", item.Title, item.Link))
	}
	return strings.Join(newsItems, "\n")
}
func sendHelpMessage(bot *tgbotapi.BotAPI, chatID int64) {
	helpMessage := "Доступные команды:\n" +
		"/help - Показать список команд с описаниями\n" +
		"/photoNews - Случайная фотоновость\n" +
		"/history - Случайный анекдот из истории\n" +
		"/anekdot - Случайный анекдот\n" +
		"/rmp - Случайный мобильный фото-пост с Reddit\n" +
		"/rmeme - Случайный мем с Reddit\n" +
		"/meme - Случайный мем\n" +
		"/tvbox - Последние новости из рубрики TV Box\n" +
		"/breaking - Последние актуальные новости\n" +
		"/auto - Последние новости из мира автомобилей\n" +
		"/bash- Последние новости из bashOrg\n" +
		"/gadgets - Последние новости из мира гаджетов"

	sendMessage(bot, chatID, helpMessage)
}
