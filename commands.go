package main

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "play",
		Description: "Plays a song",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "identifier",
				Description: "The song link or search query",
				Required:    true,
			},
		},
	},
	{
		Name:        "pause",
		Description: "Pauses/unpauses the current song",
	},
	{
		Name:        "resume",
		Description: "Pauses/unpauses the current song",
	},
	{
		Name:        "now-playing",
		Description: "Shows the current playing song",
	},
	{
		Name:        "stop",
		Description: "Stops the current song and stops the player",
	},
	{
		Name:        "skip",
		Description: "Skip the current song",
	},
	{
		Name:        "players",
		Description: "Shows all active players",
	},
	{
		Name:        "shuffle",
		Description: "Shuffles the current queue",
	},
	{
		Name:        "queue",
		Description: "Shows the current queue",
	},
	{
		Name:        "clear-queue",
		Description: "Clears the current queue",
	},
	{
		Name:        "queue-type",
		Description: "Sets the queue type",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "The queue type",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "default",
						Value: "default",
					},
					{
						Name:  "repeat-track",
						Value: "repeat-track",
					},
					{
						Name:  "repeat-queue",
						Value: "repeat-queue",
					},
				},
			},
		},
	},
	{
		Name:        "shutdown",
		Description: "Triggers the bot to shutdown, hopefully to be automatically restarted.",
	},
}

func registerCommands(s *discordgo.Session) {
	if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildId, commands); err != nil {
		slog.Error("error while registering commands", slog.Any("err", err))
	}
}
