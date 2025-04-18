package commands

import "github.com/bwmarrin/discordgo"

func init() {
	RegisterCommand("jawa", JawaCommand)
}

func JawaCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Keluar Jawa")
}
