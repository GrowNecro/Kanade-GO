package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Tipe function standar untuk command
type CommandFunc func(s *discordgo.Session, m *discordgo.MessageCreate)

// Registry semua command
var CommandRegistry = map[string]CommandFunc{}

func RegisterCommand(name string, handler CommandFunc) {
	CommandRegistry[name] = handler
}

var CommandMap = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
	"ingfo": IngfoCommand,
}
