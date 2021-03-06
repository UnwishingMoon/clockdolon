package bot

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/UnwishingMoon/clockdolon/pkg/app"
	"github.com/UnwishingMoon/clockdolon/pkg/cetus"
	"github.com/UnwishingMoon/clockdolon/pkg/db"
	"github.com/bwmarrin/discordgo"
)

var tk *time.Ticker
var dg *discordgo.Session

// Start is used to start the discord handler and bot
func Start() {
	var err error

	// Creating the bot
	dg, err = discordgo.New("Bot " + app.Conf.Bot.Token)
	if err != nil {
		log.Fatalf("[FATAL] Error during bot creation: %s", err.Error())
	}

	// Handler for messages
	dg.AddHandler(MessageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildPresences | discordgo.IntentsGuilds

	// Starting the connection to discord
	if err = dg.Open(); err != nil {
		log.Fatalf("[FATAL] Error during bot initialization: %s", err.Error())
	}

	tk = time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-tk.C:
				alertScheduled()
			}
		}
	}()
}

// Close is invoked to shutdown all sockets and connections
func Close() {
	tk.Stop()
	dg.Close()
}

// MessageCreate is used to handle all messages received from discord
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignores other bots message and itself
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	// Prefix has to be the one set
	if len(m.Content) == 0 || string(m.Content[0]) != app.Conf.Bot.Prefix {
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
	case "link":
		linkCommand(s, m)
	}
}

func linkCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	roles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "Something went wrong from our end. Please try again later!"))
		return
	}

	for _, role := range roles {
		for _, mrole := range m.Member.Roles {
			if role.ID == mrole {
				if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator ||
					role.Permissions&discordgo.PermissionManageServer == discordgo.PermissionManageServer {
					// It has the permission to execute this command
					err := db.LinkChannel(m.GuildID, m.ChannelID)
					if err != nil {
						s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "Something went wrong from our end. Please try again later!"))
						return
					}

					s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "**Channel linked**!\nAlerts will be posted here!"))
					return
				}
			}
		}
	}

	s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "You need the **Manage Server** or **Administrator** permission to do this!"))
	return
}

func timeCommand(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) {
	var description string
	durationPassed := time.Duration(time.Since(cetus.World.DayStart)/cetus.FullDay) * cetus.FullDay

	if math.Mod(time.Since(cetus.World.DayStart).Seconds(), 150*60) < 100*60 {
		// Day
		remaining := time.Until(cetus.World.NightStart.Add(durationPassed)).Truncate(1 * time.Second)
		description = fmt.Sprintf("`%s` remaining until **night**!", remaining)
	} else {
		// Night
		remaining := time.Until(cetus.World.NightEnd.Add(durationPassed)).Truncate(1 * time.Second)
		description = fmt.Sprintf("`%s` remaining until the **end** of the **night**!", remaining)
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: description,
			Color:       8359053,
		},
		Reference: m.Reference(),
	})
}

func alertCommand(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) {
	description := "You will be notified `%v minutes` before **night**!"

	if len(cmd) < 2 {
		s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "Sorry. The command needs an argument (1-60)"))
		return
	}

	if !db.GuildIsLinked(m.GuildID) {
		s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "Sorry. You have to **link** a channel before using this command, use `!help` for all the available commands"))
		return
	}

	minutes, err := strconv.Atoi(cmd[1])
	if err != nil || minutes < 1 || minutes > 60 {
		s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "Sorry. The argument specified is invalid. Only numbers from 1 to 60 are allowed."))
		return
	}

	if db.UserAlertExist(m.GuildID, m.Author.ID) {
		s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "Sorry. You already have another alert set. Remove that one before trying again."))
		return
	}

	err = db.AddUserAlert(m.GuildID, m.Author.ID, minutes)
	if err != nil {
		s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "Something went wrong from our end. Please try again later!"))
		return
	}

	if pr, err := s.State.Presence(m.GuildID, m.Author.ID); err == nil && pr.Status == discordgo.StatusOffline {
		description += "\n\n**You have to be online to receive a notification from the alert!**"
	}

	s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, description, minutes))
}

func removeCommand(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) {
	db.RemoveUserAlert(m.GuildID, m.Author.ID)

	s.ChannelMessageSendComplex(m.ChannelID, sendMessage(m, "Previous alert removed!"))
}

func helpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	const description = `
	**Clockdolon** is a simple bot that keep track of **Warframe Cetus Time** and warns you when night is about to happen.

	**Commands**
	` + "`!help`" + ` to print this message
	` + "`!time`" + ` to print the time until night
	` + "`!alert`" + ` followed by the time you want to be alerted (1-60 minutes)
	` + "`!remove`" + ` to remove yourself from the alert
	` + "`!link`" + ` to receive alerts on the channel, it must be set to be able to use alerts!!

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

func sendMessage(m *discordgo.MessageCreate, description string, args ...interface{}) *discordgo.MessageSend {
	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: fmt.Sprintf(description, args...),
			Color:       8359053,
		},
		Reference: m.Reference(),
	}

	return ms
}

func alertScheduled() {
	var timeStr string
	minutes := cetus.WorldTime()

	// Skip if out of my interval
	if minutes < 1 || minutes > 60 {
		return
	}

	if minutes == 1 {
		timeStr = "minute"
	} else {
		timeStr = "minutes"
	}

	rows, err := db.ScheduledAlerts(minutes)
	if err != nil {
		log.Printf("[Warn] Could not scan rows: %s", err.Error())
		return
	}
	defer rows.Close()

	onlineUsers := make(map[string][]string)

	for rows.Next() {
		var (
			users   string
			channel string
			guild   string
		)

		if err := rows.Scan(&users, &channel, &guild); err != nil {
			log.Printf("[Warn] Could not scan rows: %s", err.Error())
			continue // Should not use it
		}

		for _, v := range strings.Split(users, ",") {
			prs, err := dg.State.Presence(guild, v)
			if err != nil {
				log.Printf("[Warn] Could not get user presence: %s", err.Error())
				continue
			}

			if prs.Status != discordgo.StatusOffline {
				if _, prs := onlineUsers[channel]; !prs {
					onlineUsers[channel] = make([]string, 0)
				}

				onlineUsers[channel] = append(onlineUsers[channel], v)
			}
		}
	}

	for channel, u := range onlineUsers {
		if len(u) > 0 {
			description := fmt.Sprintf("`%v %s` before the **night**!\n\n<@%s>", minutes, timeStr, strings.Join(u, "> <@"))

			dg.ChannelMessageSend(channel, description)
		}
	}
}
