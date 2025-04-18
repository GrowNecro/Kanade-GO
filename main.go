package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"mydiscordbot/commands" // Ganti sesuai nama module kamu
)

var prefix string

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == commands.BotID || m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(args) == 0 {
		return
	}

	cmd := strings.ToLower(args[0])
	if handler, ok := commands.CommandRegistry[cmd]; ok {
		handler(s, m)
	} else {
		s.ChannelMessageSend(m.ChannelID, "❌ Perintah tidak ditemukan.")
	}
}

func main() {
	// Load .env
	_ = godotenv.Load()
	Token := os.Getenv("DISCORD_TOKEN")
	prefix = os.Getenv("BOT_PREFIX")

	if Token == "" || prefix == "" {
		fmt.Println("DISCORD_TOKEN atau BOT_PREFIX belum diset di .env")
		return
	}

	// Buat session Discord
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Gagal membuat session,", err)
		return
	}

	// Tangkap event Ready untuk simpan BotID
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		commands.BotID = r.User.ID
		fmt.Println("Bot aktif sebagai:", r.User.Username, "ID:", r.User.ID)
	})

	// Tambahkan handler pesan
	dg.AddHandler(handleMessage)

	// Buka koneksi
	err = dg.Open()
	if err != nil {
		fmt.Println("Gagal membuka koneksi:", err)
		return
	}

	fmt.Println("Bot sedang berjalan. Tekan CTRL+C untuk keluar.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	_ = dg.Close()
}
