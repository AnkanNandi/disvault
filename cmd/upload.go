package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AnkanNandi/disvault/app"
	"github.com/AnkanNandi/disvault/core"
	"github.com/AnkanNandi/disvault/db"
	"github.com/spf13/cobra"
)

// Flag variables
var (
	inputFile string
	groupID   int
	groupName string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file by splitting it into chunks and registering it in the database",
	Long: `This command splits a large file into chunks, uploads each chunk, and registers them in the database.
You must create groups before assigning them with flags`,
	Run: runUploadCmd,
}

// Init function to define flags and add the command to the root
func init() {
	uploadCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Path to the input file to be uploaded (required)")
	uploadCmd.Flags().IntVarP(&groupID, "id", "i", 1, "Group ID for the upload, defaults to 1 which is `uncategorized`")
	uploadCmd.Flags().StringVarP(&groupName, "name", "n", "uncategorized", "Group Name for the upload file, defaults to `uncategorized`")

	// Only one of the flags can be chosen
	uploadCmd.MarkFlagRequired("file")
	uploadCmd.MarkFlagsMutuallyExclusive("id", "name")

	// Add the upload command to the root command
	rootCmd.AddCommand(uploadCmd)
}

// runUploadCmd executes the upload command logic
func runUploadCmd(cmd *cobra.Command, args []string) {
	db.InitDatabase()
	app.Init()
	// If the group name is provided, find the corresponding group ID
	if cmd.Flags().Changed("name") {
		groupID = fetchGroupName(groupName)
	}

	// Validate the provided group ID
	ValidateGroupID(groupID)

	// Call the Upload function from the core package
	if err := core.Upload(inputFile, groupID); err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
	fmt.Println("File uploaded successfully.")
}

// fetchGroupIDByName retrieves the group ID based on the group name from the database.
func fetchGroupName(gName string) int {
	var gID int
	err := db.DB.QueryRow("SELECT group_id FROM groups WHERE group_name = ?", gName).Scan(&gID)
	switch {
	case err == sql.ErrNoRows:
		log.Fatalf("No group found with name: %s", gName)
	case err != nil:
		log.Fatalf("Failed to query group ID: %v", err)
	}
	return gID
}

// validateGroupID checks if the provided group ID exists in the database.
func ValidateGroupID(gid int) {
	err := db.DB.QueryRow("SELECT group_id FROM groups WHERE group_id = ?", gid).Scan(&gid)
	switch {
	case err == sql.ErrNoRows:
		log.Fatalf("No group found with ID: %d", gid)
	case err != nil:
		log.Fatalf("Failed to query group ID: %v", err)
	}
}
