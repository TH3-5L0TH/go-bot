package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"

	"github.com/disgoorg/disgolink/v3/disgolink"
)

var (
	urlPattern    *regexp.Regexp
	searchPattern *regexp.Regexp
	Token         string
	GuildId       string
	NodeName      string
	NodeAddress   string
	NodePassword  string
	NodeSecure    bool
)

func init() {
	fmt.Println("INITIALISING")
	godotenv.Load()

	urlPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)
	Token = os.Getenv("BOT_TOKEN")
	GuildId = os.Getenv("GUILD_ID")
	NodeName = os.Getenv("NODE_NAME")
	NodeAddress = os.Getenv("NODE_ADDRESS")
	NodePassword = os.Getenv("NODE_PASSWORD")
	NodeSecure, _ = strconv.ParseBool(os.Getenv("NODE_SECURE"))
}

type Bot struct {
	Session  *discordgo.Session
	Lavalink disgolink.Client
	Handlers map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error
	Queues   *QueueManager
}

func main() {
	slog.Info("starting discordgo example...")
	slog.Info("discordgo version:", slog.String("version", discordgo.VERSION))
	slog.Info("disgolink version:", slog.String("version", disgolink.Version))

	b := &Bot{
		Queues: &QueueManager{
			queues: make(map[string]*Queue),
		},
	}

	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		slog.Error("error while creating discordgo session", slog.Any("err", err))
		os.Exit(1)
	}
	b.Session = session

	session.State.TrackVoice = true
	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentsGuildVoiceStates

	session.AddHandler(b.onApplicationCommand)
	session.AddHandler(b.onVoiceStateUpdate)
	session.AddHandler(b.onVoiceServerUpdate)

	if err = session.Open(); err != nil {
		slog.Error("error while opening session", slog.Any("err", err))
		os.Exit(1)
	}
	defer session.Close()

	registerCommands(session)

	b.Lavalink = disgolink.New(snowflake.MustParse(session.State.User.ID),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
	)
	b.Handlers = map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error{
		"play":        b.play,
		"pause":       b.pause,
		"now-playing": b.nowPlaying,
		"stop":        b.stop,
		"queue":       b.queue,
		"clear-queue": b.clearQueue,
		"queue-type":  b.queueType,
		"shuffle":     b.shuffle,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     NodeName,
		Address:  NodeAddress,
		Password: NodePassword,
		Secure:   NodeSecure,
	})
	if err != nil {
		slog.Error("failed to add node", slog.Any("err", err))
		os.Exit(1)
	}
	version, err := node.Version(ctx)
	if err != nil {
		slog.Error("failed to get node version", slog.Any("err", err))
		os.Exit(1)
	}
	slog.Info("node version:", slog.String("version", version))

	slog.Info("DiscordGo example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func (b *Bot) onApplicationCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	data := event.ApplicationCommandData()

	handler, ok := b.Handlers[data.Name]
	if !ok {
		slog.Info("unknown command", slog.String("command", data.Name))
		return
	}
	if err := handler(event, data); err != nil {
		slog.Error("error handling command: ", slog.Any("err", err))
	}
}

func (b *Bot) onVoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if event.UserID != session.State.User.ID {
		return
	}

	var channelID *snowflake.ID
	if event.ChannelID != "" {
		id := snowflake.MustParse(event.ChannelID)
		channelID = &id
	}
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), snowflake.MustParse(event.GuildID), channelID, event.SessionID)
	if event.ChannelID == "" {
		b.Queues.Delete(event.GuildID)
	}
}

func (b *Bot) onVoiceServerUpdate(session *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), snowflake.MustParse(event.GuildID), event.Token, event.Endpoint)
}
