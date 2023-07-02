package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"

	"math/rand"
	"time"
)

const (
	chance       = "chance"
	commandStart = "start"
	info_yes_no  = "rand"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart: // если бот узнал команду, например /start, меняем текст сообщения и выводим.
		return b.handleStartCommand(message)
	case info_yes_no:
		return b.handle_yes_no(message)
	case chance:
		return b.handle_randomizer(message)
	//case zapros:
	//return b.question_for_gpt(message)
	default: //если бот не знает команды, по дефолту выводим
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidURL
	}
	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnAuthorized
	}
	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		return errUnableToSave
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.LinkSaved)
	_, err = b.bot.Send(msg)
	return err
}
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg) // _, переменная, которую мы никогда не используем
	return err
}

func (b *Bot) handle_yes_no(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	rand.Seed(time.Now().UnixNano())
	args := msg.Text[5:]
	answer := []string{"да", "нет"}[rand.Intn(2)]
	msg_ans := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintln("Ответ на вопрос: ", args, "–", answer))
	msg_ans.ReplyToMessageID = message.MessageID

	_, err := b.bot.Send(msg_ans)
	return err
}

func (b *Bot) handle_randomizer(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	rand.Seed(time.Now().UnixNano())
	args := msg.Text[7:]
	msg_ans := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Вероятность, что%s - %d%%\n", args, rand.Intn(101)))
	msg_ans.ReplyToMessageID = message.MessageID
	_, err := b.bot.Send(msg_ans)
	return err
}

//func (b *Bot) question_for_gpt(message *tgbotapi.Message) error {
//	client := openai.NewClient("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")
//	resp, err := client.CreateChatCompletion(
//		context.Background(),
//		openai.ChatCompletionRequest{
//			Model:       openai.GPT3Dot5Turbo0301,
//			Temperature: 1.0,
//			MaxTokens:   10000,
//			Messages: []openai.ChatCompletionMessage{
//				{
//					Role:    openai.ChatMessageRoleAssistant,
//					Content: message.Text[5:],
//				},
//			},
//		},
//	)
//
//	if err != nil {
//		return err
//	}
//	if len(resp.Choices) > 0 {
//		msg := tgbotapi.NewMessage(message.Chat.ID, resp.Choices[0].Message.Content)
//		msg.ReplyToMessageID = message.MessageID
//		_, err = b.bot.Send(msg)
//		if err != nil {
//			log.Fatal(err)
//		}
//	} else {
//		log.Fatal("OpenAI response is empty")
//	}
//	return err
//}
