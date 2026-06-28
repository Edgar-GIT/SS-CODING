package mainbot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var botUserID string

func registerHandlers(session *discordgo.Session) {
	session.AddHandler(onReady)
	session.AddHandler(onGuildMemberAdd)
	session.AddHandler(onMessageCreate)
	session.AddHandler(onMessageUpdate)
	session.AddHandler(onMessageDelete)
}

func onReady(session *discordgo.Session, ready *discordgo.Ready) {
	if ready.User != nil {
		botUserID = ready.User.ID
	}
	botLog("Main bot online as %s", ready.User.Username)
}

func onGuildMemberAdd(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	if event.User == nil || event.User.Bot {
		return
	}

	userID := event.User.ID
	username := event.User.Username
	if event.Member != nil && event.Member.Nick != "" {
		username = event.Member.Nick
	}

	botLog("New member: %s (%s)", username, userID)

	if err := sendWelcomeImage(session); err != nil {
		botLogError("welcome image: %v", err)
	}

	if _, err := session.ChannelMessageSend(WelcomeChannelID, buildWelcomeMessage(userID)); err != nil {
		botLogError("welcome message: %v", err)
	}
}

func onMessageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event == nil || event.Message == nil {
		return
	}
	msg := event.Message
	if msg.Author != nil && msg.Author.ID == botUserID {
		return
	}
	if msg.Author != nil && msg.Author.Bot {
		return
	}
	if msg.ChannelID == MusicChannelID {
		return
	}

	cacheStore(msg)

	if handleCommand(session, msg) {
		return
	}

	go logMessageSent(session, msg)
}

func onMessageUpdate(session *discordgo.Session, before, after *discordgo.MessageUpdate) {
	if after == nil || after.Message == nil {
		return
	}
	msg := after.Message
	if msg.Author != nil && (msg.Author.Bot || msg.Author.ID == botUserID) {
		return
	}
	if msg.ChannelID == MusicChannelID {
		return
	}

	var beforeMsg *discordgo.Message
	if before != nil {
		beforeMsg = before.Message
	}
	cacheStore(msg)
	go logMessageEdited(session, beforeMsg, msg)
}

func onMessageDelete(session *discordgo.Session, event *discordgo.MessageDelete) {
	if event == nil {
		return
	}

	snap := cacheLoad(event.ID)
	if snap == nil {
		snap = &messageSnapshot{
			MessageID: event.ID,
			GuildID:   event.GuildID,
			ChannelID: event.ChannelID,
			Content:   "_Content unavailable (not cached)_",
		}
	}

	if snap.ChannelID == MusicChannelID {
		cacheRemove(event.ID)
		return
	}
	if snap.AuthorID == botUserID {
		cacheRemove(event.ID)
		return
	}

	go logMessageDeleted(session, snap)
	cacheRemove(event.ID)
}

func sendWelcomeImage(session *discordgo.Session) error {
	path, err := welcomeImagePath()
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = session.ChannelMessageSendComplex(WelcomeChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{Name: "welcome.png", Reader: file},
		},
	})
	return err
}

func buildWelcomeMessage(userID string) string {
	mention := fmt.Sprintf("<@%s>", userID)
	rules := fmt.Sprintf("<#%s>", RulesChannelID)
	return fmt.Sprintf(
		"👋 **Ayyy wassup %s — welcome to the server, brah!** 😎\n\n"+
			"Before you dive in, swing by %s and read the rules.\n"+
			"Enjoy your stay, meet the crew, and step up your skills. Let's go! 🔥😎",
		mention, rules,
	)
}
