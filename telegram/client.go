package telegram

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"alquiler-scrapping/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Bot *tgbotapi.BotAPI
}

func NewTelegramBot() (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		return nil, err
	}

	bot.Debug = os.Getenv("DEBUG") == "true"

	return &Telegram{
		Bot: bot,
	}, nil
}

func (tg *Telegram) SendToChannel(entry database.Entry) (tgbotapi.Message, error) {
	chatId, err := strconv.ParseInt(os.Getenv("TG_CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatal("TG_CHAT_ID environment variable is incorrrect.")
	}

	text := fmt.Sprintf("%d €\n%s", entry.Price, entry.Url)

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "HTML"
	return tg.Bot.Send(msg)
}
