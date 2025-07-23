# Telegram Channel Summarizer Bot

A Telegram bot that creates humorous AI-powered summaries of channel posts.

## Features

- Monitors Telegram channel comment groups
- Automatically summarizes channel posts using OpenAI
- Limits summaries to 30% of original text length
- Creates entertaining and witty summaries in stand-up comedy style
- Adds creative ratings based on post topic

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create `.env` file from template:
   ```bash
   cp .env.example .env
   ```

4. Fill in environment variables:
   - `TELEGRAM_BOT_TOKEN` - your bot token from @BotFather
   - `OPENAI_API_KEY` - your OpenAI API key

5. Customize the AI prompt (optional):
   - Edit `PROMPT.md` to modify the summarization style
   - The bot reads the prompt template from the code block in this file

## Usage

```bash
go run main.go
```

## Bot Setup

1. Create a bot via @BotFather
2. Add the bot to your channel's comment group as an administrator
3. Grant the bot permission to read messages
4. Start the bot

The bot will automatically process new channel posts and create humorous summaries as replies to the original posts.