package commands

import "github.com/bwmarrin/discordgo"

func init() {
	RegisterCommand("hello", HelloCommand)
}

func HelloCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Hai")
}
