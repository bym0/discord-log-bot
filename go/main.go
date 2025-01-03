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

	// Add event handlers
	dg.AddHandler(voiceStateUpdate)
	dg.AddHandlerOnce(func(s *discordgo.Session, ready *discordgo.Ready) {
		// Send a message to the log channel when the bot is ready
		_, err := s.ChannelMessageSend(logChannel, "Bot is up and running!")
		if err != nil {
			fmt.Println("Error! Sending startup message:", err)
		}
	})

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

func voiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	logChannel := os.Getenv("LOGBOT_CHANNEL")
	if logChannel == "" {
		fmt.Println("Error! LOGBOT_CHANNEL is missing.")
		return
	}

	userMention := "<@" + vsu.UserID + ">"

	// Check for join (nil or empty BeforeUpdate.ChannelID -> valid ChannelID)
	if vsu.BeforeUpdate == nil || vsu.BeforeUpdate.ChannelID == "" {
		if vsu.ChannelID != "" {
			_, err := s.ChannelMessageSend(logChannel, fmt.Sprintf("%s has joined channel <#%s>.", userMention, vsu.ChannelID))
			if err != nil {
				fmt.Println("Error! Sending join message:", err)
			}
		}
		return
	}

	// Check for leave (valid BeforeUpdate.ChannelID -> nil or empty ChannelID)
	if vsu.ChannelID == "" {
		if vsu.BeforeUpdate.ChannelID != "" {
			_, err := s.ChannelMessageSend(logChannel, fmt.Sprintf("%s has left channel <#%s>.", userMention, vsu.BeforeUpdate.ChannelID))
			if err != nil {
				fmt.Println("Error! Sending leave message:", err)
			}
		}
	}
}
