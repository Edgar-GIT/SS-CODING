package musicbot

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type pendingTrack struct {
	Track     *Track
	UserID    string
	ChannelID string
	GuildID   string
}

var (
	pendingMu sync.Mutex
	pending   = map[string]*pendingTrack{}
	botUserID string
	geniusKey string
)

func registerMusicHandlers(session *discordgo.Session, cfg MusicConfig) {
	geniusKey = cfg.GeniusToken
	session.AddHandler(onReady)
	session.AddHandler(onMessageCreate)
	session.AddHandler(onInteractionCreate)
	session.AddHandler(onVoiceStateUpdate)
}

func onReady(session *discordgo.Session, _ *discordgo.Ready) {
	botUserID = session.State.User.ID
	botLog("🤖 Music bot online as %s", session.State.User.Username)
}

func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author == nil || message.Author.Bot || !strings.HasPrefix(message.Content, commandPrefix) {
		return
	}
	if message.GuildID == "" || message.ChannelID != CommandChannelID {
		return
	}

	fields := strings.Fields(message.Content)
	if len(fields) == 0 {
		return
	}
	command := strings.ToLower(strings.TrimPrefix(fields[0], commandPrefix))
	args := strings.TrimSpace(strings.TrimPrefix(message.Content, fields[0]))
	channelID := message.ChannelID
	guildID := message.GuildID
	author := message.Author

	switch command {
	case "help":
		handleHelp(session, channelID)
	case "playmusic":
		handlePlayMusic(session, channelID, guildID, author, args)
	case "playsc":
		handlePlaySC(session, channelID, guildID, author, args, false)
	case "playscurl":
		handlePlaySC(session, channelID, guildID, author, args, true)
	case "replay":
		getPlayer(guildID).setLoop(true)
		sendPlain(session, channelID, "🔁 Loop enabled.")
	case "stoploop":
		getPlayer(guildID).setLoop(false)
		sendPlain(session, channelID, "⏹ Loop stopped.")
	case "volume":
		handleVolume(session, channelID, guildID, args)
	case "lyrics":
		handleLyrics(session, channelID, guildID, args)
	case "createplaylist":
		handleCreatePlaylist(session, channelID, author.ID, args)
	case "playlistadd":
		handlePlaylistAdd(session, channelID, author.ID, args)
	case "showplaylist":
		handleShowPlaylist(session, channelID, author.ID)
	case "playlist":
		handleQueuePlaylist(session, channelID, guildID, author)
	}
}

func handleHelp(session *discordgo.Session, channelID string) {
	embed := &discordgo.MessageEmbed{
		Title: "🎵 Music Bot Commands",
		Color: panelColor,
		Fields: []*discordgo.MessageEmbedField{
			{Name: commandPrefix + "playmusic <name>", Value: "Search and play music from YouTube", Inline: false},
			{Name: "⏯ / ⏭ / ⏹ / ⬅ / ➡", Value: "Control music with the panel buttons", Inline: false},
			{Name: commandPrefix + "replay", Value: "Loop the current song", Inline: false},
			{Name: commandPrefix + "stoploop", Value: "Stop looping", Inline: false},
			{Name: commandPrefix + "lyrics <music>", Value: "Fetch lyrics", Inline: false},
			{Name: commandPrefix + "volume <0-100>", Value: "Set bot volume", Inline: false},
			{Name: commandPrefix + "createplaylist <name>", Value: "Create your personal playlist", Inline: false},
			{Name: commandPrefix + "playlistadd <music>", Value: "Add music to your playlist", Inline: false},
			{Name: commandPrefix + "showplaylist", Value: "Show your playlist", Inline: false},
			{Name: commandPrefix + "playlist", Value: "Queue your playlist", Inline: false},
			{Name: commandPrefix + "playsc <name>", Value: "Stream from SoundCloud search", Inline: false},
			{Name: commandPrefix + "playscurl <url>", Value: "Stream a SoundCloud URL", Inline: false},
		},
	}
	_, _ = session.ChannelMessageSendEmbed(channelID, embed)
}

func handlePlayMusic(session *discordgo.Session, channelID, guildID string, author *discordgo.User, query string) {
	if query == "" {
		sendPlain(session, channelID, "⚠️ Provide a song name.")
		return
	}

	track, err := searchYouTube(query)
	if err != nil || track == nil {
		sendPlain(session, channelID, "❌ Music not found.")
		return
	}

	pendingMu.Lock()
	pending[author.ID] = &pendingTrack{Track: track, UserID: author.ID, ChannelID: channelID, GuildID: guildID}
	pendingMu.Unlock()

	_, _ = session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content:    fmt.Sprintf("<@%s>, is this the music you were looking for?", author.ID),
		Embed:      buildPreviewEmbed(track),
		Components: confirmComponents(author.ID),
	})
}

func handlePlaySC(session *discordgo.Session, channelID, guildID string, author *discordgo.User, query string, direct bool) {
	if query == "" {
		sendPlain(session, channelID, "⚠️ Provide a search or URL.")
		return
	}

	label := query
	if !direct {
		sendPlain(session, channelID, "🔎 Searching SoundCloud for: **"+query+"**")
	}
	track, err := searchSoundCloud(query)
	if err != nil || track == nil || track.URL == "" {
		sendPlain(session, channelID, "❌ Could not find this track on SoundCloud.")
		return
	}

	gp := getPlayer(guildID)
	gp.enqueue(*track)
	sendPlain(session, channelID, "➕ Added to the queue (SoundCloud): **"+label+"**")
	startPlayback(session, channelID, guildID, author)
}

func handleVolume(session *discordgo.Session, channelID, guildID, args string) {
	value, err := strconv.Atoi(strings.TrimSpace(args))
	if err != nil || value < 0 || value > 100 {
		sendPlain(session, channelID, "⚠️ Volume must be between 0 and 100.")
		return
	}
	gp := getPlayer(guildID)
	gp.setVolume(float64(value) / 100)
	sendPlain(session, channelID, fmt.Sprintf("🔊 Volume set to %d%%", value))
}

func handleLyrics(session *discordgo.Session, channelID, guildID, args string) {
	query := strings.TrimSpace(args)
	if query == "" {
		current, _, _, _, _, _ := getPlayer(guildID).snapshot()
		if current != nil {
			query = current.Title
		}
	}
	if query == "" {
		sendPlain(session, channelID, "⚠️ Usage: "+commandPrefix+"lyrics <song name>")
		return
	}
	if geniusKey == "" {
		sendPlain(session, channelID, "❌ GENIUS_TOKEN missing in .env")
		return
	}

	chunks, source, err := fetchLyrics(query, geniusKey)
	if err != nil {
		if source != "" {
			sendPlain(session, channelID, "⚠️ Lyrics couldn't be fetched, check them here: "+source)
		} else {
			sendPlain(session, channelID, "❌ "+err.Error())
		}
		return
	}
	for _, chunk := range chunks {
		sendPlain(session, channelID, "```"+chunk+"```")
	}
	sendPlain(session, channelID, "📌 Source: "+source)
}

func handleCreatePlaylist(session *discordgo.Session, channelID, userID, name string) {
	if strings.TrimSpace(name) == "" {
		sendPlain(session, channelID, "⚠️ Example: "+commandPrefix+"createplaylist MyPlaylist")
		return
	}
	if err := createUserPlaylist(userID, name); err != nil {
		sendPlain(session, channelID, "⚠️ "+err.Error())
		return
	}
	sendPlain(session, channelID, "✅ Playlist '"+name+"' created.")
}

func handlePlaylistAdd(session *discordgo.Session, channelID, userID, song string) {
	if strings.TrimSpace(song) == "" {
		sendPlain(session, channelID, "⚠️ Provide a song name to add.")
		return
	}
	if err := addSongToPlaylist(userID, song); err != nil {
		sendPlain(session, channelID, "⚠️ "+err.Error())
		return
	}
	sendPlain(session, channelID, "➕ Added to your playlist: "+song)
}

func handleShowPlaylist(session *discordgo.Session, channelID, userID string) {
	playlist, err := loadUserPlaylist(userID)
	if err != nil || playlist == nil {
		sendPlain(session, channelID, "⚠️ You don't have a playlist.")
		return
	}
	if len(playlist.Songs) == 0 {
		embed := &discordgo.MessageEmbed{Title: "📂 " + playlist.Name, Color: panelColor}
		embed.Fields = []*discordgo.MessageEmbedField{{Name: "No songs in playlist", Value: "Add songs with " + commandPrefix + "playlistadd", Inline: false}}
		_, _ = session.ChannelMessageSendEmbed(channelID, embed)
		return
	}

	header := &discordgo.MessageEmbed{
		Title: "📂 " + playlist.Name, Color: panelColor,
		Fields: []*discordgo.MessageEmbedField{{Name: "Total songs", Value: strconv.Itoa(len(playlist.Songs)), Inline: false}},
	}
	_, _ = session.ChannelMessageSendEmbed(channelID, header)

	for i := 0; i < len(playlist.Songs); i += 10 {
		end := i + 10
		if end > len(playlist.Songs) {
			end = len(playlist.Songs)
		}
		embed := &discordgo.MessageEmbed{Color: panelColor}
		for j, song := range playlist.Songs[i:end] {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name: fmt.Sprintf("%d.", i+j+1), Value: song, Inline: false,
			})
		}
		_, _ = session.ChannelMessageSendEmbed(channelID, embed)
	}
}

func handleQueuePlaylist(session *discordgo.Session, channelID, guildID string, author *discordgo.User) {
	playlist, err := loadUserPlaylist(author.ID)
	if err != nil || playlist == nil {
		sendPlain(session, channelID, "⚠️ You don't have a playlist.")
		return
	}
	gp := getPlayer(guildID)
	for _, song := range playlist.Songs {
		gp.enqueueQuery(song)
	}
	sendPlain(session, channelID, "🎶 Playlist '"+playlist.Name+"' queued.")
	startPlayback(session, channelID, guildID, author)
}

func startPlayback(session *discordgo.Session, channelID, guildID string, _ *discordgo.User) {
	gp := getPlayer(guildID)
	if err := gp.connect(session, VoiceChannelID); err != nil {
		sendPlain(session, channelID, "❌ Could not connect to voice channel.")
		return
	}
	gp.startIfIdle(session, channelID)
	sendOrUpdatePanel(session, gp, channelID)
}

func onInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionMessageComponent {
		return
	}
	if interaction.ChannelID != CommandChannelID {
		return
	}
	data := interaction.MessageComponentData()
	customID := data.CustomID
	guildID := interaction.GuildID
	channelID := interaction.ChannelID
	userID := interaction.Member.User.ID
	gp := getPlayer(guildID)

	switch {
	case customID == "music_pause":
		handlePause(session, interaction, gp)
	case customID == "music_skip":
		gp.skip(session, channelID)
		replyEphemeral(session, interaction, "⏭ Skipped")
		refreshPanel(session, guildID)
	case customID == "music_stop":
		gp.stopAll(session, channelID)
		replyEphemeral(session, interaction, "⏹ Stopped and disconnected")
		refreshPanel(session, guildID)
	case customID == "music_loop":
		enabled := gp.toggleLoop()
		replyEphemeral(session, interaction, "🔁 Loop "+onOff(enabled))
		refreshPanel(session, guildID)
	case customID == "music_vol_down":
		value := gp.adjustVolume(-0.05)
		replyEphemeral(session, interaction, fmt.Sprintf("🔉 Volume: %d%%", int(value*100)))
		refreshPanel(session, guildID)
	case customID == "music_vol_up":
		value := gp.adjustVolume(0.05)
		replyEphemeral(session, interaction, fmt.Sprintf("🔊 Volume: %d%%", int(value*100)))
		refreshPanel(session, guildID)
	case customID == "music_prev":
		if gp.playPrevious(session, channelID) {
			replyEphemeral(session, interaction, "▶ Playing previous song.")
		} else {
			replyEphemeral(session, interaction, "🚫 No previous song.")
		}
		refreshPanel(session, guildID)
	case customID == "music_next":
		gp.mu.Lock()
		hasQueue := len(gp.queue) > 0
		gp.mu.Unlock()
		if hasQueue {
			gp.skip(session, channelID)
			replyEphemeral(session, interaction, "▶ Playing next song.")
		} else {
			replyEphemeral(session, interaction, "🚫 No more songs in the queue.")
		}
		refreshPanel(session, guildID)
	case strings.HasPrefix(customID, "music_confirm_yes_"):
		handleConfirmYes(session, interaction, userID, channelID, guildID)
	case strings.HasPrefix(customID, "music_confirm_no_"):
		handleConfirmNo(session, interaction, userID)
	}
}

func handlePause(session *discordgo.Session, interaction *discordgo.InteractionCreate, gp *GuildPlayer) {
	gp.mu.Lock()
	vc := gp.vc
	gp.mu.Unlock()
	if vc == nil {
		replyEphemeral(session, interaction, "🚫 Nothing is playing.")
		return
	}
	if gp.togglePause() {
		replyEphemeral(session, interaction, "⏸ Paused")
	} else {
		replyEphemeral(session, interaction, "▶ Resumed")
	}
	refreshPanel(session, gp.guildID)
}

func handleConfirmYes(session *discordgo.Session, interaction *discordgo.InteractionCreate, userID, channelID, guildID string) {
	if !strings.HasSuffix(interaction.MessageComponentData().CustomID, userID) {
		replyEphemeral(session, interaction, "⚠️ Only the requester can confirm.")
		return
	}
	pendingMu.Lock()
	item := pending[userID]
	delete(pending, userID)
	pendingMu.Unlock()
	if item == nil || item.Track == nil {
		replyEphemeral(session, interaction, "⌛ Request expired.")
		return
	}

	gp := getPlayer(guildID)
	gp.enqueue(*item.Track)
	replyEphemeral(session, interaction, "➕ Added to the queue.")
	sendPlain(session, channelID, "➕ Added to the queue: **"+item.Track.Title+"**")

	member := interaction.Member
	if member != nil {
		startPlayback(session, channelID, guildID, member.User)
	}
}

func handleConfirmNo(session *discordgo.Session, interaction *discordgo.InteractionCreate, userID string) {
	if !strings.HasSuffix(interaction.MessageComponentData().CustomID, userID) {
		replyEphemeral(session, interaction, "⚠️ Only the requester can decline.")
		return
	}
	pendingMu.Lock()
	delete(pending, userID)
	pendingMu.Unlock()
	replyEphemeral(session, interaction, "❌ Music declined.")
}

func onVoiceStateUpdate(session *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.UserID == botUserID && vs.ChannelID == "" {
		return
	}
	if vs.ChannelID == "" {
		return
	}
	gp := getPlayer(vs.GuildID)
	gp.mu.Lock()
	vc := gp.vc
	voiceChannelID := gp.voiceChannelID
	gp.mu.Unlock()
	if vc == nil || voiceChannelID != vs.ChannelID {
		return
	}

	members := 0
	guild, err := session.State.Guild(vs.GuildID)
	if err != nil || guild == nil {
		return
	}
	for _, state := range guild.VoiceStates {
		if state.ChannelID == voiceChannelID && state.UserID != botUserID {
			members++
		}
	}
	if members > 0 {
		return
	}

	go func(guildID, channelID string) {
		time.Sleep(30 * time.Second)
		gp := getPlayer(guildID)
		gp.mu.Lock()
		vc := gp.vc
		voiceChannelID := gp.voiceChannelID
		gp.mu.Unlock()
		if vc == nil {
			return
		}
		alone := 0
		guild, err := session.State.Guild(guildID)
		if err != nil || guild == nil {
			return
		}
		for _, state := range guild.VoiceStates {
			if state.ChannelID == voiceChannelID && state.UserID != botUserID {
				alone++
			}
		}
		if alone == 0 {
			gp.stopAll(session, channelID)
			botLog("💤 Auto-disconnected due to inactivity.")
		}
	}(vs.GuildID, "")
}

func onOff(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}
