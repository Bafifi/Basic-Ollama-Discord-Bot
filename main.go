package main

import (
	bot "aibotlocal/DiscordBotHandler"
	ollama "aibotlocal/OllamaHandler"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	var err = godotenv.Load()

	if err != nil {
		fmt.Println(err)
	}

	var token = os.Getenv("DISCORD_TOKEN")

	modelMap := map[string]string{
		"discordbot": "./modelfiles/modelfile_discordbot",
	}
	bot.BotToken = token
	ollama.CreateModel_url = "http://localhost:11434/api/create"
	ollama.GenerateResponse_url = "http://localhost:11434/api/generate"
	ollama.CreateModels(modelMap)
	bot.RunBot()
}
