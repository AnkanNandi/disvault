package cmd

import (
	"fmt"
	"log"

	"github.com/AnkanNandi/disvault/app"
	"github.com/AnkanNandi/disvault/core"
	"github.com/AnkanNandi/disvault/db"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <file_id>",
	Short: "Delete files using their IDs",
	Long: `Delete a file using its registered ID, similar to how downloads work.

File parts are deleted first, followed by the main file registration in the database.
For example, if a file has 10 parts, all parts will be deleted before the main file registration is removed.`,
	Args: cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run:  runDeleteCmd,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

// runDeleteCmd executes the delete command logic
func runDeleteCmd(cmd *cobra.Command, args []string) {
	db.InitDatabase()
	app.Init()
	// Parse the file ID from the argument
	fileID, err := ParseFileID(args[0])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		cmd.Help()
		return
	}

	// Fetch the file to ensure it exists before deletion
	if _, err := FetchFileNameByID(fileID); err != nil {
		log.Fatalf("Error fetching file: %v", err)
	}

	// Delete the file parts and then the main file registration
	if err := core.DeleteFileParts(fileID); err != nil {
		log.Fatalf("Failed to delete file: %v", err)
	}

	fmt.Println("File deleted successfully.")
}
