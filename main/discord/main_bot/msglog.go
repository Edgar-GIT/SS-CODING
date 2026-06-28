package mainbot

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	colorLogSent   = 0x57F287
	colorLogEdit   = 0xFEE75C
	colorLogDelete = 0xED4245
)

func shouldLogMessage(msg *discordgo.Message) bool {
	if msg == nil || msg.GuildID == "" {
		return false
	}
	if msg.Author != nil && msg.Author.Bot {
		return false
	}
	if msg.ChannelID == MusicChannelID {
		return false
	}
	return true
}

func logMessageSent(session *discordgo.Session, msg *discordgo.Message) {
	if !shouldLogMessage(msg) {
		return
	}
	channelMention := fmt.Sprintf("<#%s>", msg.ChannelID)
	author := authorLabel(msg)
	content := strings.TrimSpace(msg.Content)
	if content == "" {
		content = "_No text_"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📩 Message sent",
		Color:       colorLogSent,
		Description: content,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Author", Value: author, Inline: true},
			{Name: "Channel", Value: channelMention, Inline: true},
			{Name: "Message ID", Value: msg.ID, Inline: true},
		},
		Timestamp: msg.Timestamp.Format(time.RFC3339),
		Footer:    &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Guild %s", msg.GuildID)},
	}
	if len(msg.Attachments) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Attachments",
			Value: attachmentSummary(msg.Attachments),
		})
	}

	sendLogPayload(session, embed, msg.Attachments, msg.Embeds)
}

func logMessageEdited(session *discordgo.Session, before, after *discordgo.Message) {
	if after == nil || !shouldLogMessage(after) {
		return
	}
	beforeText := "_No text_"
	if before != nil && strings.TrimSpace(before.Content) != "" {
		beforeText = before.Content
	}
	afterText := strings.TrimSpace(after.Content)
	if afterText == "" {
		afterText = "_No text_"
	}

	embed := &discordgo.MessageEmbed{
		Title: "✏️ Message edited",
		Color: colorLogEdit,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Author", Value: authorLabel(after), Inline: true},
			{Name: "Channel", Value: fmt.Sprintf("<#%s>", after.ChannelID), Inline: true},
			{Name: "Message ID", Value: after.ID, Inline: true},
			{Name: "Before", Value: truncateField(beforeText), Inline: false},
			{Name: "After", Value: truncateField(afterText), Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	var attachments []*discordgo.MessageAttachment
	if len(after.Attachments) > 0 {
		attachments = after.Attachments
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Attachments",
			Value: attachmentSummary(attachments),
		})
	}
	sendLogPayload(session, embed, attachments, after.Embeds)
}

func logMessageDeleted(session *discordgo.Session, snap *messageSnapshot) {
	if snap == nil || snap.ChannelID == MusicChannelID {
		return
	}
	if snap.AuthorID != "" && snap.AuthorID == botUserID {
		return
	}

	content := strings.TrimSpace(snap.Content)
	if content == "" {
		content = "_No text_"
	}
	author := snap.AuthorTag
	if snap.AuthorID != "" {
		author = fmt.Sprintf("%s (`%s`)", snap.AuthorTag, snap.AuthorID)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🗑️ Message deleted",
		Color:       colorLogDelete,
		Description: truncateField(content),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Author", Value: author, Inline: true},
			{Name: "Channel", Value: fmt.Sprintf("<#%s>", snap.ChannelID), Inline: true},
			{Name: "Message ID", Value: snap.MessageID, Inline: true},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	if len(snap.Attachments) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Attachments",
			Value: attachmentSummary(snap.Attachments),
		})
	}
	sendLogPayload(session, embed, snap.Attachments, snap.Embeds)
}

func sendLogPayload(session *discordgo.Session, embed *discordgo.MessageEmbed, attachments []*discordgo.MessageAttachment, embeds []*discordgo.MessageEmbed) {
	if _, err := session.ChannelMessageSendEmbed(LogChannelID, embed); err != nil {
		botLogError("log embed: %v", err)
		return
	}
	for _, e := range embeds {
		if e == nil {
			continue
		}
		copyEmbed := *e
		if _, err := session.ChannelMessageSendEmbed(LogChannelID, &copyEmbed); err != nil {
			botLogError("log forwarded embed: %v", err)
		}
	}
	for _, att := range attachments {
		sendLogAttachment(session, att)
	}
}

func sendLogAttachment(session *discordgo.Session, att *discordgo.MessageAttachment) {
	if att == nil || att.URL == "" {
		return
	}
	if strings.HasPrefix(att.ContentType, "image/") {
		_, _ = session.ChannelMessageSendEmbed(LogChannelID, &discordgo.MessageEmbed{
			Title: "📎 Image attachment",
			Image: &discordgo.MessageEmbedImage{URL: att.URL},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%s · %s", att.Filename, formatSize(att.Size)),
			},
		})
		return
	}

	resp, err := http.Get(att.URL)
	if err != nil {
		_, _ = session.ChannelMessageSend(LogChannelID, fmt.Sprintf("📎 Attachment: %s\n%s", att.Filename, att.URL))
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		_, _ = session.ChannelMessageSend(LogChannelID, fmt.Sprintf("📎 Attachment: %s\n%s", att.Filename, att.URL))
		return
	}

	_, err = session.ChannelMessageSendComplex(LogChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("📎 **%s** (%s)", att.Filename, formatSize(att.Size)),
		Files: []*discordgo.File{{
			Name:        att.Filename,
			ContentType: att.ContentType,
			Reader:      bytes.NewReader(data),
		}},
	})
	if err != nil {
		_, _ = session.ChannelMessageSend(LogChannelID, fmt.Sprintf("📎 Attachment: %s\n%s", att.Filename, att.URL))
	}
}

func authorLabel(msg *discordgo.Message) string {
	if msg == nil || msg.Author == nil {
		return "Unknown"
	}
	return fmt.Sprintf("%s (`%s`)", msg.Author.String(), msg.Author.ID)
}

func attachmentSummary(attachments []*discordgo.MessageAttachment) string {
	if len(attachments) == 0 {
		return "_None_"
	}
	var b strings.Builder
	for i, att := range attachments {
		if att == nil {
			continue
		}
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("• **%s** (%s)", att.Filename, formatSize(att.Size)))
	}
	return b.String()
}

func truncateField(text string) string {
	const max = 1000
	if len(text) <= max {
		return text
	}
	return text[:max] + "…"
}

func formatSize(size int) string {
	if size <= 0 {
		return "unknown size"
	}
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
}
