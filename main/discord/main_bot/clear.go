package mainbot

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func handleClear(session *discordgo.Session, msg *discordgo.Message) {
	if msg == nil || msg.Author == nil {
		return
	}
	cfg, err := LoadMainConfig()
	if err != nil {
		return
	}
	if cfg.OwnerID == "" {
		sendTemp(session, msg.ChannelID, "⚠️ `MAIN_OWNER_ID` is not set in `.env`.")
		return
	}
	if msg.Author.ID != cfg.OwnerID {
		sendTemp(session, msg.ChannelID, "❌ You can't use this command.")
		return
	}

	_ = session.ChannelMessageDelete(msg.ChannelID, msg.ID)

	deleted := 0
	for round := 0; round < 50; round++ {
		messages, err := session.ChannelMessages(msg.ChannelID, 100, "", "", "")
		if err != nil || len(messages) == 0 {
			break
		}

		ids := make([]string, 0, len(messages))
		for _, m := range messages {
			ids = append(ids, m.ID)
		}

		if err := session.ChannelMessagesBulkDelete(msg.ChannelID, ids); err != nil {
			for _, id := range ids {
				if err := session.ChannelMessageDelete(msg.ChannelID, id); err == nil {
					deleted++
				}
				time.Sleep(350 * time.Millisecond)
			}
		} else {
			deleted += len(ids)
		}

		if len(messages) < 100 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	sendTemp(session, msg.ChannelID, fmt.Sprintf("🗑️ %d messages deleted.", deleted))
	botLog("!clear in channel %s — %d messages removed", msg.ChannelID, deleted)
}

func sendTemp(session *discordgo.Session, channelID, content string) {
	reply, err := session.ChannelMessageSend(channelID, content)
	if err != nil || reply == nil {
		return
	}
	time.AfterFunc(5*time.Second, func() {
		_ = session.ChannelMessageDelete(channelID, reply.ID)
	})
}

func handleCommand(session *discordgo.Session, msg *discordgo.Message) bool {
	if msg == nil || !strings.HasPrefix(msg.Content, commandPrefix) {
		return false
	}
	fields := strings.Fields(msg.Content)
	if len(fields) == 0 {
		return false
	}
	cmd := strings.ToLower(fields[0])
	switch cmd {
	case commandPrefix + "clear":
		handleClear(session, msg)
		return true
	default:
		return false
	}
}
