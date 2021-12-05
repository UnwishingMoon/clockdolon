package bot

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/UnwishingMoon/clockdolon/pkg/app"
	"github.com/UnwishingMoon/clockdolon/pkg/cetus"
	"github.com/bwmarrin/discordgo"
)

func Start() (*discordgo.Session, error) {
	// Creating the bot
	dg, err := discordgo.New("Bot " + app.Conf.Bot.Token)
	if err != nil {
		return nil, err
	}

	// Handler for messages
	dg.AddHandler(MessageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Starting the connection to discord
	if err = dg.Open(); err != nil {
		return nil, err
	}

	return dg, nil
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var isAdmin bool = false

	// Ignores other bots message and itself
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	// Prefix has to be the one set
	if string(m.Content[0]) != app.Conf.Bot.Prefix {
		return
	}

	// Checking if user is an admin
	for _, v := range app.Conf.Bot.Admins {
		if m.Author.ID == v {
			isAdmin = true
		}
	}

	// Temporary block
	if !isAdmin {
		return
	}

	// Removing the prefix
	cmd := strings.Fields(m.Content[len(app.Conf.Bot.Prefix):])

	// First parameter of the command
	switch cmd[0] {
	case "help":
		helpCommand(s, m)
	case "time":
		timeCommand(s, m, cmd)
	case "alert":
		alertCommand(s, m, cmd)
	case "remove":
		removeCommand(s, m, cmd)
	case "me":
		meCommand(s, m)
	}
}

func timeCommand(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) {
	var description string
	daysPassed := time.Duration(time.Since(cetus.Cetus.DayStart).Seconds() / (150 * 60) * float64(time.Second))

	if math.Mod(time.Since(cetus.Cetus.DayStart).Seconds(), 150*60) < 100*60 {
		// Day
		remaining := time.Until(cetus.Cetus.NightStart.Add(daysPassed)).Round(1 * time.Second)
		description = fmt.Sprintf("`%s`"+" remaining until **night**!", remaining)
	} else {
		// Night
		remaining := time.Until(cetus.Cetus.NightEnd.Add(daysPassed)).Round(1 * time.Second)
		description = fmt.Sprintf("`%s`"+" remaining until the end of the **night**!", remaining)
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: description,
		Color:       8359053,
	})
}

func alertCommand(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) {
	//var description string

	if len(cmd) < 2 {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Description: "Sorry. The command needs an argument (1m-60m are allowed)",
			Color:       8359053,
		})
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: "Sorry. The command needs an argument (1m-60m are allowed)",
		Color:       8359053,
	})
}

func removeCommand(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) {

}

func helpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	const description = `
	**Clockdolon** is a simple bot that keep track of **Warframe Cetus Time** and warns you when night is about to happen.

	**Commands**
	` + "`!help`" + ` to print this message
	` + "`!time`" + ` to print the time until night
	` + "`!alert`" + ` followed by the time you want to be alerted (1m-60m are allowed)
	` + "`!remove`" + ` to remove yourself from the alert

	**Support**
	If you want to help and keep the bot running, you can [donate](https://streamlabs.com/unwishingmoon/) here.
	`

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: description,
		Color:       8359053,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Clockdolon",
			IconURL: "https://www.diegocastagna.com/assets/img/projects/clockdolon-icon.bf37ry4.png",
		},
	})
}

func meCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	const description = ``

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: description,
		Color:       8359053,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Clockdolon",
			IconURL: "https://www.diegocastagna.com/assets/img/projects/clockdolon-icon.bf37ry4.png",
		},
	})
}
