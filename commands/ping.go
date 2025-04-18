package commands

import "github.com/bwmarrin/discordgo"

func init() {
	RegisterCommand("ping", PingCommand)
}

func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong! 🏓")
}
