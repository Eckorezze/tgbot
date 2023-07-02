package telegram

import (
	"awesomeProject1/pkg/config"
	"awesomeProject1/pkg/repos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repos.TokenRepository
	redirectURL     string

	messages config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tr repos.TokenRepository, redirectURL string, messages config.Messages) *Bot { //конструктор, который принимает ссылку на объект бот и присваивать в поле
	return &Bot{bot: bot, pocketClient: pocketClient, redirectURL: redirectURL, tokenRepository: tr, messages: messages}
}

func (b *Bot) Start() error { //метод
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) { //метод
	for update := range updates { //получаем сообщения из канала
		if update.Message == nil { //если наше обновление не является пользовательским сообщением, то скип
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}
		//b.handleMessage(update.Message)
		//b.question_for_gpt(update.Message)
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) { //метод
	u := tgbotapi.NewUpdate(0) //создаем новую конфигурацию для получения обновлений | long-polling
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u), nil
}
