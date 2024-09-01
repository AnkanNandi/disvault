package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/AnkanNandi/disvault/db"
	"github.com/bwmarrin/discordgo"
)

// Bot parameters
type config struct {
	BotToken  string `json:"bot_token"`
	ChannelID string `json:"channel_id"`
}

var (
	Config  config
	Session *discordgo.Session // Reuse a single session
)

// Initialize the discord bot configs
func Init() {
	// Load configuration from file
	if err := loadConfig("data/config.json"); err != nil {
		log.Fatalf("Failed to load config: %v\n\nRun `disvault setup` to add the bot token and channel id", err)
	}

	// Create a new Discord session
	var err error
	Session, err = discordgo.New("Bot " + Config.BotToken)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	// Make sure to close the session when the application stops
	defer Session.Close()
}

// Loads the config values of discord bot
func loadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	if err := json.Unmarshal(byteValue, &Config); err != nil {
		return fmt.Errorf("could not unmarshal config file: %w", err)
	}

	return nil
}

// UploadFile uploads a file chunk using an existing Discord session
func UploadFile(ctx context.Context, partName string, partsStruct *db.Parts) (string, error) {
	f, err := os.Open(partName)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	ms := &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   partName,
				Reader: f,
			},
		},
	}
	msgSent, err := Session.ChannelMessageSendComplex(Config.ChannelID, ms)
	if err != nil {
		return "", fmt.Errorf("error sending message: %w", err)
	}
	partsStruct.Parts = append(partsStruct.Parts, msgSent.ID)
	fmt.Printf("Message ID: %v\n", msgSent.ID)
	return msgSent.ID, nil
}

/*
Downloads the parts in same order as uploaded and joins each part on the fly,
since each file is added to the file before and then closed, no big memory usage is observed

	`Tested with 775mb file, hashes match`
	'TODO: Add a way to check for sudden stops in download or stich if possible i.e. program crash in middle of function'
*/
func DownloadPart(partID string) ([]byte, error) {
	msg, err := Session.ChannelMessage(Config.ChannelID, partID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving message: %w", err)
	}

	if len(msg.Attachments) == 0 {
		return nil, fmt.Errorf("no attachments found in message: %s", partID)
	}

	res, err := http.Get(msg.Attachments[0].URL)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return data, nil
}
