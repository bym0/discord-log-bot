package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var botID string

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

	// Get bot user ID for later checks
	dg.AddHandlerOnce(func(s *discordgo.Session, ready *discordgo.Ready) {
		botID = s.State.User.ID
		// Send startup message
		_, err := s.ChannelMessageSend(logChannel, "Bot is up and running!")
		if err != nil {
			fmt.Println("Error! Sending startup message:", err)
		}
	})

	// Add handlers for events
	dg.AddHandler(voiceStateUpdate)
	dg.AddHandler(handleInteraction)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error! Opening Discord connection:", err)
		return
	}

	// Register application commands
	registerCommands(dg)

	fmt.Println("Success! Bot is now running.")

	// Graceful shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up commands on shutdown
	dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, "", nil)
	dg.Close()
}

func registerCommands(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with Pong!",
		},
	}

	// Register commands globally
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			fmt.Printf("Error creating command %s: %v\n", cmd.Name, err)
		} else {
			fmt.Printf("Command %s registered.\n", cmd.Name)
		}
	}
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "ping":
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
			if err != nil {
				fmt.Println("Error responding to /ping command:", err)
			}
		}
	}
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
