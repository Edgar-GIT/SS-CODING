package mainbot

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type messageSnapshot struct {
	MessageID   string
	GuildID     string
	ChannelID   string
	ChannelName string
	AuthorID    string
	AuthorTag   string
	Content     string
	Attachments []*discordgo.MessageAttachment
	Embeds      []*discordgo.MessageEmbed
}

var (
	cacheMu sync.RWMutex
	cache   = map[string]*messageSnapshot{}
)

const maxCacheSize = 8000

func cacheStore(msg *discordgo.Message) {
	if msg == nil {
		return
	}
	snap := snapshotFromMessage(msg)
	cacheMu.Lock()
	if len(cache) >= maxCacheSize {
		for k := range cache {
			delete(cache, k)
			break
		}
	}
	cache[msg.ID] = snap
	cacheMu.Unlock()
}

func cacheLoad(messageID string) *messageSnapshot {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	return cache[messageID]
}

func cacheRemove(messageID string) {
	cacheMu.Lock()
	delete(cache, messageID)
	cacheMu.Unlock()
}

func snapshotFromMessage(msg *discordgo.Message) *messageSnapshot {
	if msg == nil {
		return nil
	}
	snap := &messageSnapshot{
		MessageID: msg.ID,
		GuildID:   msg.GuildID,
		ChannelID: msg.ChannelID,
		Content:   msg.Content,
	}
	if msg.Author != nil {
		snap.AuthorID = msg.Author.ID
		snap.AuthorTag = msg.Author.String()
	}
	if msg.ChannelID != "" {
		snap.ChannelName = channelLabel(msg)
	}
	if len(msg.Attachments) > 0 {
		snap.Attachments = append([]*discordgo.MessageAttachment(nil), msg.Attachments...)
	}
	if len(msg.Embeds) > 0 {
		snap.Embeds = append([]*discordgo.MessageEmbed(nil), msg.Embeds...)
	}
	return snap
}

func channelLabel(msg *discordgo.Message) string {
	if msg == nil {
		return "unknown"
	}
	return "#" + msg.ChannelID
}
