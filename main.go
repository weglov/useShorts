package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type Bot struct {
	tg     *tgbotapi.BotAPI
	ai     *openai.Client
	chatID int64
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	openaiToken := os.Getenv("OPENAI_API_KEY")
	if openaiToken == "" {
		log.Fatal("OPENAI_API_KEY not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Error creating bot:", err)
	}

	aiClient := openai.NewClient(openaiToken)

	app := &Bot{
		tg: bot,
		ai: aiClient,
	}

	log.Printf("Bot started: %s", bot.Self.UserName)
	log.Printf("Bot ready to receive messages...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			app.handleMessage(update.Message)
		}
		
		if update.ChannelPost != nil {
			app.handleChannelPost(update.ChannelPost)
		}
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	// Check for automatic forward from channel
	if msg.IsAutomaticForward && msg.ForwardFromChat != nil && msg.ForwardFromChat.IsChannel() {
		b.handleChannelPost(msg)
		return
	}
}

func (b *Bot) handleChannelPost(post *tgbotapi.Message) {
	if post.Text == "" {
		return
	}

	textLength := len(post.Text)
	if textLength < 100 {
		return
	}

	summary, err := b.summarizeText(post.Text)
	if err != nil {
		log.Printf("Summarization error: %v", err)
		return
	}

	response := summary

	msg := tgbotapi.NewMessage(post.Chat.ID, response)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = post.MessageID // Reply to original post

	_, err = b.tg.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *Bot) loadPromptTemplate() (string, error) {
	content, err := ioutil.ReadFile("PROMPT.md")
	if err != nil {
		return "", fmt.Errorf("failed to read PROMPT.md: %v", err)
	}

	return strings.TrimSpace(string(content)), nil
}

func (b *Bot) summarizeText(text string) (string, error) {
	maxLength := int(float64(len(text)) * 0.3)
	
	promptTemplate, err := b.loadPromptTemplate()
	if err != nil {
		return "", fmt.Errorf("failed to load prompt template: %v", err)
	}
	
	prompt := fmt.Sprintf(promptTemplate, maxLength, text)

	resp, err := b.ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: maxLength / 2,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty AI response")
	}

	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	return summary, nil
}