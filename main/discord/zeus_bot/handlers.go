package zeusbot

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	cfgMu          sync.RWMutex
	zeusCfg        ZeusConfig
	zeusActivated  bool
	activatedMu    sync.Mutex
)

func isZeus(authorID string) bool {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return authorID == zeusCfg.ZeusID
}

func setActivated(value bool) {
	activatedMu.Lock()
	zeusActivated = value
	activatedMu.Unlock()
}

func zeusIsActivated() bool {
	activatedMu.Lock()
	defer activatedMu.Unlock()
	return zeusActivated
}

func registerHandlers(session *discordgo.Session) {
	session.AddHandler(onReady)
	session.AddHandler(onMessageCreate)
}

func onReady(session *discordgo.Session, ready *discordgo.Ready) {
	botLog("⚡ ZEUS ⚡ online as %s", ready.User.Username)
}

func onMessageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event == nil || event.Message == nil || event.Author == nil || event.Author.Bot {
		return
	}
	if event.GuildID == "" || !strings.HasPrefix(event.Content, commandPrefix) {
		return
	}

	fields := strings.Fields(event.Content)
	if len(fields) == 0 {
		return
	}
	cmd := strings.ToLower(fields[0])
	args := fields[1:]

	switch cmd {
	case commandPrefix + "zeus":
		handleZeus(session, event.Message)
	case commandPrefix + "zeusbolt":
		handleZeusBolt(session, event.Message, args)
	case commandPrefix + "zeuslight":
		handleZeusLight(session, event.Message)
	case commandPrefix + "zeusban":
		handleZeusBan(session, event.Message, args)
	case commandPrefix + "help":
		handleZeusHelp(session, event.Message)
	}
}

func handleZeus(session *discordgo.Session, msg *discordgo.Message) {
	authorMention := fmt.Sprintf("<@%s>", msg.Author.ID)
	if !isZeus(msg.Author.ID) {
		_, _ = session.ChannelMessageSend(msg.ChannelID,
			fmt.Sprintf("I AM ⚡ZEUS ⚡, GOD OF LIGHT AND DESTRUCTION, AND YOU %s SHALL BE PUNISHED FOR GIVING ORDERS TO A GOD!", authorMention))
		punishUser(session, msg.GuildID, msg.Author.ID, 150)
		return
	}

	setActivated(true)
	_ = sendImage(session, msg.ChannelID, "zeusintro.png")
	_, _ = session.ChannelMessageSend(msg.ChannelID,
		"WHO DARES TO BOTHER ME? I AM ⚡ZEUS ⚡, GOD OF LIGHT AND DESTRUCTION AND I ONLY SHALL OBEY TO ORDERS FROM HERCULES !")
}

func handleZeusBolt(session *discordgo.Session, msg *discordgo.Message, args []string) {
	if !requireActivated(session, msg) {
		return
	}
	if !punishIfNotZeus(session, msg) {
		return
	}
	if len(args) == 0 {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "⚠️ You must specify a member to punish. Example: `!zeusbolt @user 60`")
		return
	}

	targetID := parseUserID(args[0])
	seconds := 60
	if len(args) > 1 {
		seconds = parsePositiveInt(args[1], 60)
	}

	member, err := session.GuildMember(msg.GuildID, targetID)
	if err != nil || member == nil {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "❌ Member not found!")
		return
	}

	_, _ = session.ChannelMessageSend(msg.ChannelID, "I AM ⚡ZEUS ⚡, AND I SHALL PUNISH YOU FROM YOUR SINS!!")
	punishUser(session, msg.GuildID, member.User.ID, seconds)
}

func handleZeusLight(session *discordgo.Session, msg *discordgo.Message) {
	if !requireActivated(session, msg) {
		return
	}
	if !punishIfNotZeus(session, msg) {
		return
	}

	_, _ = session.ChannelMessageSend(msg.ChannelID, "I AM ⚡ZEUS ⚡, AND I NOW SHALL ENLIGHTEN YOU PEASANTS WITH THE DIVINE LIGHT OF GOD!!")
	_, _ = session.ChannelMessageSend(msg.ChannelID, "BEHOLD THE DIVINE LIGHT OF GOD!!")
	_ = sendImage(session, msg.ChannelID, "zeuslight.png")
	_, _ = session.ChannelMessageSend(msg.ChannelID, "MAY THE DIVINE LIGHT GUIDE YOU TO SALVATION!!")
}

func handleZeusBan(session *discordgo.Session, msg *discordgo.Message, args []string) {
	if !requireActivated(session, msg) {
		return
	}
	if !punishIfNotZeus(session, msg) {
		return
	}
	if len(args) == 0 {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "⚠️ Usage: `!zeusban <user_id>`")
		return
	}

	targetID := parseUserID(args[0])
	if targetID == "" {
		_, _ = session.ChannelMessageSend(msg.ChannelID, "⚠️ Invalid user ID.")
		return
	}

	go func(guildID, channelID, userID string) {
		if err := sendImage(session, channelID, "ban.png"); err != nil {
			botLogError("zeusban image: %v", err)
		}
		time.Sleep(3 * time.Second)
		if err := session.GuildBanCreateWithReason(guildID, userID, "Struck down by Zeus!", 0); err != nil {
			botLogError("zeusban: %v", err)
			_, _ = session.ChannelMessageSend(channelID, fmt.Sprintf("❌ Could not ban <@%s>: %v", userID, err))
			return
		}
		_, _ = session.ChannelMessageSend(channelID, fmt.Sprintf("⚡ <@%s> has been banished by Zeus!", userID))
	}(msg.GuildID, msg.ChannelID, targetID)
}

func handleZeusHelp(session *discordgo.Session, msg *discordgo.Message) {
	if !requireActivated(session, msg) {
		return
	}
	if !punishIfNotZeus(session, msg) {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "⚡ ZEUS COMMANDS ⚡",
		Color: 0xFFD700,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "!zeus", Value: "Zeus announces his presence.", Inline: false},
			{Name: "!zeusbolt @user {time}", Value: "Silences a member for {time} seconds.", Inline: false},
			{Name: "!zeuslight", Value: "Sends the divine light.", Inline: false},
			{Name: "!zeusban <user_id>", Value: "Banishes a member from the server.", Inline: false},
			{Name: "!help", Value: "Shows this help message.", Inline: false},
		},
	}
	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, embed)
}

func requireActivated(session *discordgo.Session, msg *discordgo.Message) bool {
	if zeusIsActivated() {
		return true
	}
	_, _ = session.ChannelMessageSend(msg.ChannelID, "⚠️ Zeus is not in action yet! Use `!zeus` first.")
	return false
}

func punishIfNotZeus(session *discordgo.Session, msg *discordgo.Message) bool {
	if isZeus(msg.Author.ID) {
		return true
	}
	authorMention := fmt.Sprintf("<@%s>", msg.Author.ID)
	_, _ = session.ChannelMessageSend(msg.ChannelID,
		fmt.Sprintf("I AM ⚡ZEUS ⚡, AND YOU %s SHALL BE PUNISHED!", authorMention))
	punishUser(session, msg.GuildID, msg.Author.ID, 150)
	return false
}

func punishUser(session *discordgo.Session, guildID, userID string, seconds int) {
	if seconds <= 0 {
		seconds = 150
	}
	until := time.Now().UTC().Add(time.Duration(seconds) * time.Second)
	if err := session.GuildMemberTimeout(guildID, userID, &until); err != nil {
		botLogError("timeout %s: %v", userID, err)
		return
	}

	dm, err := session.UserChannelCreate(userID)
	if err != nil {
		return
	}
	_, _ = session.ChannelMessageSend(dm.ID, fmt.Sprintf("⚡ You have been silenced by Zeus for %d seconds!", seconds))
}

func sendImage(session *discordgo.Session, channelID, filename string) error {
	path, err := imagePath(filename)
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Files: []*discordgo.File{{Name: filename, Reader: file}},
	})
	return err
}
