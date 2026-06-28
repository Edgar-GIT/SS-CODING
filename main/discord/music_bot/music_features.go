package musicbot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (gp *GuildPlayer) setTextChannel(channelID string) {
	gp.mu.Lock()
	gp.textChannelID = channelID
	gp.mu.Unlock()
}

func (gp *GuildPlayer) setStayInChannel(enabled bool) {
	gp.mu.Lock()
	gp.stayInChannel = enabled
	gp.mu.Unlock()
}

func (gp *GuildPlayer) stayInChannelEnabled() bool {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	return gp.stayInChannel
}

func (gp *GuildPlayer) resetSkipVotes() {
	gp.mu.Lock()
	gp.skipVoters = nil
	gp.mu.Unlock()
}

func (gp *GuildPlayer) shuffleQueue() int {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	n := len(gp.queue)
	if n < 2 {
		return n
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		gp.queue[i], gp.queue[j] = gp.queue[j], gp.queue[i]
	}
	return n
}

func (gp *GuildPlayer) replayLast(session *discordgo.Session, channelID string) bool {
	gp.mu.Lock()
	if len(gp.history) == 0 {
		gp.mu.Unlock()
		return false
	}
	last := gp.history[len(gp.history)-1]
	gp.queue = append([]Track{last}, gp.queue...)
	playing := gp.playing
	gp.mu.Unlock()

	if playing {
		gp.stopPlaybackNow()
		return true
	}
	if err := gp.connect(session, VoiceChannelID); err != nil {
		return false
	}
	gp.startIfIdle(session, channelID)
	return true
}

func countVoiceMembers(session *discordgo.Session, guildID, channelID string) int {
	if channelID == "" {
		return 0
	}
	guild, err := session.State.Guild(guildID)
	if err != nil || guild == nil {
		return 0
	}
	count := 0
	for _, state := range guild.VoiceStates {
		if state.ChannelID == channelID && state.UserID != botUserID {
			count++
		}
	}
	return count
}

func isUserInVoiceChannel(session *discordgo.Session, guildID, channelID, userID string) bool {
	if channelID == "" {
		return false
	}
	guild, err := session.State.Guild(guildID)
	if err != nil || guild == nil {
		return false
	}
	for _, state := range guild.VoiceStates {
		if state.UserID == userID && state.ChannelID == channelID {
			return true
		}
	}
	return false
}

func votesRequired(memberCount int) int {
	if memberCount <= 0 {
		return 1
	}
	return (memberCount + 1) / 2
}

// voteSkip registers a vote. Returns current votes, required votes, and whether skip triggered.
func (gp *GuildPlayer) voteSkip(session *discordgo.Session, channelID, userID string) (ok bool, current, required int, skipped bool, msg string) {
	gp.mu.Lock()
	voiceChannelID := gp.voiceChannelID
	playing := gp.playing
	gp.mu.Unlock()

	if !playing || voiceChannelID == "" {
		return false, 0, 0, false, "Nothing is playing."
	}
	if !isUserInVoiceChannel(session, gp.guildID, voiceChannelID, userID) {
		return false, 0, 0, false, "You must be in the music voice channel to vote."
	}

	members := countVoiceMembers(session, gp.guildID, voiceChannelID)
	required = votesRequired(members)

	gp.mu.Lock()
	if gp.skipVoters == nil {
		gp.skipVoters = make(map[string]struct{})
	}
	if _, voted := gp.skipVoters[userID]; voted {
		current = len(gp.skipVoters)
		gp.mu.Unlock()
		return true, current, required, false, ""
	}
	gp.skipVoters[userID] = struct{}{}
	current = len(gp.skipVoters)
	gp.mu.Unlock()

	if current >= required {
		gp.resetSkipVotes()
		gp.skip(session, channelID)
		return true, current, required, true, ""
	}
	return true, current, required, false, ""
}

func (gp *GuildPlayer) reconnectVoice(session *discordgo.Session, voiceChannelID, textChannelID string) {
	gp.mu.Lock()
	if gp.reconnecting || !gp.wantVoice {
		gp.mu.Unlock()
		return
	}
	gp.reconnecting = true
	gp.mu.Unlock()

	defer func() {
		gp.mu.Lock()
		gp.reconnecting = false
		gp.mu.Unlock()
	}()

	for attempt := 1; attempt <= 3; attempt++ {
		botLogWarn("Voice reconnect attempt %d/3...", attempt)
		if err := gp.connect(session, voiceChannelID); err != nil {
			botLogError("Reconnect failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		botLogInfo("Voice reconnected to channel %s", voiceChannelID)
		if textChannelID != "" {
			sendPlain(session, textChannelID, "🔌 Reconnected to voice channel.")
		}
		return
	}
	botLogError("Voice reconnect gave up after 3 attempts")
	if textChannelID != "" {
		sendPlain(session, textChannelID, "❌ Lost voice connection and could not reconnect.")
	}
}

func buildQueueEmbed(gp *GuildPlayer) *discordgo.MessageEmbed {
	current, queue, _, looping, volume, _ := gp.snapshot()
	embed := &discordgo.MessageEmbed{
		Title: "📋 Queue",
		Color: panelColor,
	}

	if current != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Now Playing",
			Value: "**" + current.Title + "**",
		})
	}

	if len(queue) == 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Up Next", Value: "_Empty_", Inline: false,
		})
	} else {
		var lines []string
		for i, track := range queue {
			title := track.Title
			if title == "" {
				title = "Unknown"
			}
			lines = append(lines, fmt.Sprintf("**%d.** %s — %s", i+1, title, formatDuration(track.Duration)))
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("Up Next (%d)", len(queue)),
			Value: stringsJoinLines(lines, 15),
		})
	}

	embed.Fields = append(embed.Fields,
		&discordgo.MessageEmbedField{Name: "Loop", Value: loopStatus(looping), Inline: true},
		&discordgo.MessageEmbedField{Name: "Volume", Value: fmt.Sprintf("%d%%", int(volume*100)), Inline: true},
		&discordgo.MessageEmbedField{Name: "Stay in channel", Value: onOff(gp.stayInChannelEnabled()), Inline: true},
	)
	return embed
}

func stringsJoinLines(lines []string, max int) string {
	if len(lines) > max {
		extra := len(lines) - max
		lines = lines[:max]
		lines = append(lines, fmt.Sprintf("_…and %d more_", extra))
	}
	return strings.Join(lines, "\n")
}
