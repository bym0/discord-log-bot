package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discordToken := os.Getenv("LOGBOT_TOKEN")
	if discordToken == "" {
		fmt.Println("Error! LOGBOT_TOKEN is missing.")
		return
	}
	logChannel := os.Getenv("LOGBOT_CHANNEL")
	if logChannel == "" {
		fmt.Println("Error! LOGBOT_CHANNEL is missing.")
		return
	}

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Error! Creating Discord session:", err)
		return
	}

	dg.AddHandler(voiceStateUpdate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error! Opening Discord connection:", err)
		return
	}

	fmt.Println("Success! Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func voiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	logChannel := os.Getenv("LOGBOT_CHANNEL")
	if logChannel == "" {
		fmt.Println("Error! LOGBOT_CHANNEL is missing.")
		return
	}

	userName := m.Member.Mention()
	currentChannelID := m.ChannelID
	beforeUpdate := m.BeforeUpdate

	if beforeUpdate == nil {
		_, _ = s.ChannelMessageSend(logChannel, userName+" has joined channel <#"+currentChannelID+">.")
	} else {
		_, _ = s.ChannelMessageSend(logChannel, userName+" has left channel <#"+beforeUpdate.ChannelID+">.")
	}
}
