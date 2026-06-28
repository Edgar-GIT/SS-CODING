package musicbot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const panelColor = 0x57F287

func buildPanelEmbed(gp *GuildPlayer) *discordgo.MessageEmbed {
	current, queue, _, looping, volume, page := gp.snapshot()
	embed := &discordgo.MessageEmbed{
		Title: "🎵 Music Bot",
		Color: panelColor,
	}

	if current != nil {
		pos, dur := gp.progressSnapshot()
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Now Playing",
			Value:  fmt.Sprintf("**%s**", current.Title),
			Inline: false,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Progress",
			Value:  buildProgressBar(pos, float64(dur)),
			Inline: false,
		})
		if current.Thumbnail != "" {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: current.Thumbnail}
		}
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Now Playing", Value: "No music playing", Inline: false,
		})
	}

	if len(queue) > 0 {
		start := page * 10
		end := start + 10
		if end > len(queue) {
			end = len(queue)
		}
		totalQueued := len(queue)
		for i, song := range queue[start:end] {
			title := song.Title
			if title == "" {
				title = song.Query
			}
			if title == "" {
				title = "Unknown"
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%d. %s", start+i+1, title),
				Value:  fmt.Sprintf("Duration: %s", formatDuration(song.Duration)),
				Inline: false,
			})
		}
		if totalQueued > 10 {
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Queue page %d/%d · %d up next", page+1, gp.pageCount(totalQueued), totalQueued),
			}
		} else if current != nil {
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%d up next after current song", totalQueued),
			}
		}
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Queue", Value: "Empty", Inline: false,
		})
	}

	embed.Fields = append(embed.Fields,
		&discordgo.MessageEmbedField{Name: "🔊 Volume", Value: fmt.Sprintf("%d%%", int(volume*100)), Inline: true},
		&discordgo.MessageEmbedField{Name: "🔁 Loop", Value: loopStatus(looping), Inline: true},
		&discordgo.MessageEmbedField{Name: "📌 Stay", Value: onOff(gp.stayInChannelEnabled()), Inline: true},
	)
	return embed
}

func onOff(enabled bool) string {
	if enabled {
		return "On ✅"
	}
	return "Off ❌"
}

func loopStatus(enabled bool) string {
	if enabled {
		return "Enabled ✅"
	}
	return "Disabled ❌"
}

func buildProgressBar(elapsed, total float64) string {
	const width = 18
	if total <= 0 {
		return fmt.Sprintf("`%s` %s", strings.Repeat("─", width), formatTimestamp(elapsed))
	}
	if elapsed < 0 {
		elapsed = 0
	}
	if elapsed > total {
		elapsed = total
	}
	ratio := elapsed / total
	pos := int(ratio * float64(width-1))
	if pos < 0 {
		pos = 0
	}

	var b strings.Builder
	b.WriteString("`")
	for i := 0; i < width; i++ {
		switch {
		case i < pos:
			b.WriteRune('▬')
		case i == pos:
			b.WriteRune('🔘')
		default:
			b.WriteRune('─')
		}
	}
	b.WriteString("`\n")
	b.WriteString(formatTimestamp(elapsed))
	b.WriteString(" / ")
	b.WriteString(formatTimestamp(total))
	return b.String()
}

func formatTimestamp(sec float64) string {
	if sec < 0 {
		sec = 0
	}
	m := int(sec) / 60
	s := int(sec) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

func panelComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			button("music_pause", "⏯ Pause/Resume", discordgo.SuccessButton),
			button("music_skip", "⏭ Skip", discordgo.PrimaryButton),
			button("music_stop", "⏹ Stop", discordgo.DangerButton),
			button("music_loop", "🔁 Loop", discordgo.SecondaryButton),
			button("music_stay", "📌 Stay", discordgo.SecondaryButton),
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			button("music_rewind", "⏪ -5s", discordgo.SecondaryButton),
			button("music_forward", "⏩ +5s", discordgo.SecondaryButton),
			button("music_vol_down", "🔉 -5%", discordgo.SecondaryButton),
			button("music_vol_up", "🔊 +5%", discordgo.SecondaryButton),
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			button("music_prev", "⬅ Prev Song", discordgo.SecondaryButton),
			button("music_next", "Next Song ➡", discordgo.SecondaryButton),
			button("music_replay_last", "🔂 Replay Last", discordgo.SecondaryButton),
			button("music_lyrics", "📝 Lyrics", discordgo.SecondaryButton),
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			button("music_vote_skip", "🗳️ Vote Skip", discordgo.PrimaryButton),
			button("music_shuffle", "🔀 Shuffle", discordgo.SecondaryButton),
			button("music_clear_queue", "🗑️ Clear Queue", discordgo.DangerButton),
			button("music_view_queue", "📋 View Queue", discordgo.SecondaryButton),
		}},
	}
}

func button(id, label string, style discordgo.ButtonStyle) discordgo.Button {
	return discordgo.Button{CustomID: id, Label: label, Style: style}
}

func sendOrUpdatePanel(session *discordgo.Session, gp *GuildPlayer, channelID string) {
	embed := buildPanelEmbed(gp)
	components := panelComponents()
	panelChannel, panelMessage := gp.panelRef()

	if panelChannel != "" && panelMessage != "" {
		_, err := session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel:    panelChannel,
			ID:         panelMessage,
			Embed:      embed,
			Components: &components,
		})
		if err == nil {
			return
		}
	}

	msg, err := session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed:      embed,
		Components: components,
	})
	if err == nil {
		gp.setPanel(msg.ChannelID, msg.ID)
	}
}

func refreshPanel(session *discordgo.Session, guildID string) {
	gp := getPlayer(guildID)
	channelID, messageID := gp.panelRef()
	if channelID == "" || messageID == "" {
		return
	}
	embed := buildPanelEmbed(gp)
	components := panelComponents()
	_, _ = session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: channelID, ID: messageID, Embed: embed, Components: &components,
	})
}

func buildPreviewEmbed(track *Track) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: track.Title,
		Color: panelColor,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Uploader", Value: fallback(track.Uploader, "N/A"), Inline: true},
			{Name: "Duration", Value: formatDuration(track.Duration), Inline: true},
		},
	}
	if track.Thumbnail != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: track.Thumbnail}
	}
	return embed
}

func confirmComponents(userID string) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			button("music_confirm_yes_"+userID, "Yes", discordgo.SuccessButton),
			button("music_confirm_no_"+userID, "No", discordgo.DangerButton),
		}},
	}
}

func fallback(value, alt string) string {
	if strings.TrimSpace(value) == "" {
		return alt
	}
	return value
}

func sendPlain(session *discordgo.Session, channelID, content string) {
	_, _ = session.ChannelMessageSend(channelID, content)
}

func replyEphemeral(session *discordgo.Session, interaction *discordgo.InteractionCreate, content string) {
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func replyEphemeralEmbed(session *discordgo.Session, interaction *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}
