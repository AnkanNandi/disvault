package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/AnkanNandi/disvault/app"
	"github.com/AnkanNandi/disvault/core"
	"github.com/AnkanNandi/disvault/db"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <file_id>",
	Short: "Download files using their IDs",
	Long: `Download command allows for downloading files using their ID only.
The files are saved in the 'output' folder with the same name as during upload.

Example usage:
	disvault download <file_id>`,
	Args: cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run:  runDownloadCmd,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

// runDownloadCmd executes the download command logic
func runDownloadCmd(cmd *cobra.Command, args []string) {
	db.InitDatabase()
	app.Init()
	// Convert the file ID argument from string to integer
	fileID, err := ParseFileID(args[0])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		cmd.Help()
		return
	}

	// Check if the file exists and retrieve its name
	fileName, err := FetchFileNameByID(fileID)
	if err != nil {
		log.Fatalf("Failed to fetch file: %v", err)
	}

	// Download and reassemble the file
	if err := core.DownloadAndReassembleFile(fileID, fileName); err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}

	fmt.Println("File downloaded successfully.")
}

// parseFileID converts a string file ID to an integer and validates it
func ParseFileID(fileIDStr string) (int, error) {
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid file ID '%s'. It must be an integer", fileIDStr)
	}
	return fileID, nil
}

// fetchFileNameByID checks if the file exists in the database and returns its name
func FetchFileNameByID(id int) (string, error) {
	var fileName string
	err := db.DB.QueryRow("SELECT name FROM files WHERE id = ?", id).Scan(&fileName)
	switch {
	case err == sql.ErrNoRows:
		return "", fmt.Errorf("no file found with ID: %d", id)
	case err != nil:
		return "", fmt.Errorf("error fetching file name: %v", err)
	}
	return fileName, nil
}
