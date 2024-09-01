package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AnkanNandi/disvault/db"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

// Flags for the setup command
var (
	botToken  string
	channelID string
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the bot configuration",
	Long:  `Setup command allows you to configure the bot token and the channel ID for file outputs.`,
	Run:   runSetupCmd,
}

func init() {
	setupCmd.Flags().StringVarP(&botToken, "token", "t", "", "Discord bot token")
	setupCmd.Flags().StringVarP(&channelID, "channel", "c", "", "Discord channel ID")
	setupCmd.MarkFlagRequired("token")   // Make the token flag mandatory
	setupCmd.MarkFlagRequired("channel") // Make the channel flag mandatory

	rootCmd.AddCommand(setupCmd)
}

func runSetupCmd(cmd *cobra.Command, args []string) {
	db.InitDatabase()
	// Ensure both flags are provided
	if botToken == "" || channelID == "" {
		fmt.Println("Error: Both -t (token) and -c (channel) flags are required.")
		return
	}

	// Initialize Discord session
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Open a connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord: %v", err)
	}
	defer dg.Close()

	// Test sending a message to the specified channel
	testMessage := "Bot setup test message."
	_, err = dg.ChannelMessageSend(channelID, testMessage)
	if err != nil {
		log.Fatalf("Error sending message to channel: %v", err)
	}

	// Prepare configuration data
	config := map[string]string{
		"bot_token":  botToken,
		"channel_id": channelID,
	}

	// Create the data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	// Write configuration to config.json
	configPath := filepath.Join(dataDir, "config.json")
	configFile, err := os.Create(configPath)
	if err != nil {
		log.Fatalf("Error creating config.json file: %v", err)
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(config); err != nil {
		log.Fatalf("Error writing configuration to config.json: %v", err)
	}

	fmt.Printf("Configuration saved successfully in %s.\n", configPath)
}
