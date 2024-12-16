package main

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) shutdownNo(event *discordgo.InteractionCreate, data discordgo.MessageComponentInteractionData) error {
	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Shutdown canceled",
		},
	})
}

func (b *Bot) shutdownYes(event *discordgo.InteractionCreate, data discordgo.MessageComponentInteractionData) error {
	b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Bot shutting down",
		},
	})
	
	b.exit()
	return nil
}