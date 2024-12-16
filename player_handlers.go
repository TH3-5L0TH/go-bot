package main

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Bot) onPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	slog.Info("PlayerPause:", slog.Any("Event", event))
	//fmt.Printf("onPlayerPause: %v\n", event)
}

func (b *Bot) onPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	slog.Info("PlayerResume:", slog.Any("Event", event))
	//fmt.Printf("onPlayerResume: %v\n", event)
}

func (b *Bot) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	slog.Info("TrackStart:", slog.Any("Event", event))
	//fmt.Printf("onTrackStart: %v\n", event)
}

func (b *Bot) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	slog.Info("TrackEnd:", slog.Any("Event", event))
	//fmt.Printf("onTrackEnd: %v\n", event)

	if !event.Reason.MayStartNext() {
		return
	}

	queue := b.Queues.Get(event.GuildID().String())
	var (
		nextTrack lavalink.Track
		ok        bool
	)
	switch queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case QueueTypeRepeatTrack:
		nextTrack = event.Track

	case QueueTypeRepeatQueue:
		queue.Add(event.Track)
		nextTrack, ok = queue.Next()
	}

	if !ok {
		return
	}
	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		slog.Error("Failed to play next track", slog.Any("err", err))
	}
}

func (b *Bot) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	slog.Info("TrackException:", slog.Any("Event", event))
	//slog.Info("TrackException", slog.String("Cause", *event.Exception.Cause), slog.String("Message", *&event.Exception.Message))
	//fmt.Printf("onTrackException: %v\n", event)
}

func (b *Bot) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	slog.Info("TrackStuck:", slog.Any("Event", event))
	//fmt.Printf("onTrackStuck: %v\n", event)
}

func (b *Bot) onWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	slog.Info("WebSocketClosed:", slog.Any("Event", event))
	//fmt.Printf("onWebSocketClosed: %v\n", event)
}
