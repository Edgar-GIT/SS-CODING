package welcomebot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func registerHandlers(session *discordgo.Session) {
	session.AddHandler(onReady)
	session.AddHandler(onGuildMemberAdd)
}

func onReady(session *discordgo.Session, ready *discordgo.Ready) {
	botLog("Welcome bot online as %s", ready.User.Username)
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
