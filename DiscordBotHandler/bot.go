package bot

import (
	ollama "aibotlocal/OllamaHandler"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		fmt.Printf("Discord Bot Error: %v", e)
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.SessionID {
		return
	}

	switch {
	case strings.Contains(m.Content, "!ping"):
		s.ChannelMessageSend(m.ChannelID, "pong!")
	case strings.Contains(m.Content, "!discordbot"):
		prompt := strings.ReplaceAll(m.Content, "!discordbot", "")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🤔")
		response := ollama.GenerateResponse(prompt, "discordbot")
		s.ChannelMessageSend(m.ChannelID, response)
		s.MessageReactionRemove(m.ChannelID, m.ID, "🤔", s.State.User.ID)
	}
}

func RunBot() {
	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(messageHandler)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running until there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
