package commands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// BotID akan di-set saat event Ready diterima
var BotID string

func IngfoCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("[DEBUG] IngfoCommand dipanggil")

	// Parsing argument
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Gunakan command `!ingfo atas` untuk melihat BC atas")
		return
	}
	location := strings.Join(args[1:], " ")
	fmt.Println("[DEBUG] Location:", location)

	// Ambil API URL dan Token dari environment
	apiURL := os.Getenv("CCTV_API_URL")
	apiToken := os.Getenv("CCTV_API_TOKEN")
	if apiURL == "" || apiToken == "" {
		s.ChannelMessageSend(m.ChannelID, "API CCTV belum di-setting.")
		return
	}

	// Menambahkan reaction untuk menunjukkan proses sedang berjalan
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, "🔄"); err != nil {
		fmt.Println("Gagal menambahkan reaction 🔄:", err)
	}

	// Menyiapkan URL untuk request
	url := fmt.Sprintf("%singfo/%s", apiURL, location)
	fmt.Println("[DEBUG] URL yang digunakan untuk API:", url)

	// Membuat request ke API
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error membuat request:", err)
		s.ChannelMessageSend(m.ChannelID, "Terjadi kesalahan saat membuat permintaan ke server.")
		handleFailReaction(s, m)
		return
	}
	req.Header.Set("Authorization", apiToken)

	// Mengirim request dan menerima response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error saat request:", err)
		handleFailReaction(s, m)
		s.ChannelMessageSend(m.ChannelID, "Gagal menghubungi server CCTV.")
		return
	}
	defer resp.Body.Close()

	// Mengecek jika status code bukan 200
	if resp.StatusCode != 200 {
		fmt.Println("Status code dari server:", resp.StatusCode)
		handleFailReaction(s, m)
		s.ChannelMessageSend(m.ChannelID, "Gagal mendapatkan data dari server.")
		return
	}

	// Menangani response JSON
	var rawResult struct {
		Count json.RawMessage `json:"count"`
		Image string          `json:"image"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rawResult); err != nil {
		fmt.Println("Error decoding JSON:", err)
		handleFailReaction(s, m)
		s.ChannelMessageSend(m.ChannelID, "Gagal membaca data dari server.")
		return
	}

	// Parsing count ke int
	var count int
	if err := json.Unmarshal(rawResult.Count, &count); err != nil {
		var countStr string
		if err := json.Unmarshal(rawResult.Count, &countStr); err != nil {
			fmt.Println("Gagal parsing count sebagai string:", err)
			handleFailReaction(s, m)
			s.ChannelMessageSend(m.ChannelID, "Data count dari server tidak bisa dipahami.")
			return
		}
		if countStr == "-" || countStr == "" {
			fmt.Println("[DEBUG] Count berisi tanda '-' atau kosong, dianggap 0")
			count = 0
		} else {
			count, err = strconv.Atoi(countStr)
			if err != nil {
				fmt.Println("Gagal konversi count string ke int:", err)
				handleFailReaction(s, m)
				s.ChannelMessageSend(m.ChannelID, "Data count tidak valid.")
				return
			}
		}
	}
	fmt.Println("[DEBUG] Detected count:", count)

	// Decode gambar base64
	imageBytes, err := base64.StdEncoding.DecodeString(rawResult.Image)
	if err != nil {
		fmt.Println("Error decode base64 image:", err)
		handleFailReaction(s, m)
		s.ChannelMessageSend(m.ChannelID, "Gagal memproses gambar dari server.")
		return
	}

	// Membuat file untuk attachment
	file := &discordgo.File{
		Name:        fmt.Sprintf("cctv_%s.jpg", location),
		ContentType: "image/jpeg",
		Reader:      bytes.NewReader(imageBytes),
	}

	// Kirim pesan dengan file
	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("BC %s: Detected %d person(s).", location, count),
		Files:   []*discordgo.File{file},
	})
	if err != nil {
		fmt.Println("Error kirim pesan dengan gambar:", err)
		handleFailReaction(s, m)
		return
	}

	// Hapus reaction 🔄 dan tambahkan 👍
	if err := s.MessageReactionRemove(m.ChannelID, m.ID, "🔄", s.State.User.ID); err != nil {
		fmt.Println("Gagal hapus reaction 🔄:", err)
	}

	if err := s.MessageReactionAdd(m.ChannelID, m.ID, "👍"); err != nil {
		fmt.Println("Gagal tambahkan reaction 👍:", err)
	}

}

// Fungsi untuk menambahkan reaction error (❌)
func handleFailReaction(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Remove the reaction safely, handling any errors gracefully
	if err := s.MessageReactionRemove(m.ChannelID, m.ID, "🔄", s.State.User.ID); err != nil {
		fmt.Println("Gagal hapus reaction 🔄:", err)
	}

	// Try adding the error reaction (❌)
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, "❌"); err != nil {
		fmt.Println("Gagal tambahkan reaction ❌:", err)
	}
}

// Inisialisasi command
func init() {
	RegisterCommand("ingfo", IngfoCommand)
}
