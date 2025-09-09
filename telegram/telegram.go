package telegram

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	chatID  string
	enabled bool
}

func NewBot(token, chatID string) (*Bot, error) {
	if token == "" || chatID == "" {
		log.Println("Telegram bot not configured - alerts will be logged only")
		return &Bot{enabled: false}, nil
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	log.Printf("Telegram bot authorized on account %s", bot.Self.UserName)

	return &Bot{
		bot:     bot,
		chatID:  chatID,
		enabled: true,
	}, nil
}

func (b *Bot) SendMessage(message string) error {
	if !b.enabled {
		log.Printf("TELEGRAM ALERT (not sent - bot disabled): %s", message)
		return nil
	}

	chatID, err := strconv.ParseInt(b.chatID, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chat ID: %w", err)
	}

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = false

	if _, err := b.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send Telegram message: %w", err)
	}

	log.Printf("Telegram alert sent successfully")
	return nil
}

func (b *Bot) SendPriceAlert(productName, platform string, oldPrice, newPrice float64, url string) error {
	message := fmt.Sprintf(
		"🚨 <b>PRICE DROP ALERT!</b> 🚨\n\n"+
			"📦 <b>Product:</b> %s\n"+
			"🏪 <b>Platform:</b> %s\n"+
			"💰 <b>Previous Price:</b> ₹%.2f\n"+
			"💸 <b>Current Price:</b> ₹%.2f\n"+
			"💵 <b>Savings:</b> ₹%.2f\n\n"+
			"🔗 <a href=\"%s\">View Product</a>",
		productName, platform, oldPrice, newPrice, oldPrice-newPrice, url,
	)

	return b.SendMessage(message)
}

func (b *Bot) IsEnabled() bool {
	return b.enabled
}
