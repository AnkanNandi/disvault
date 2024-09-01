package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/AnkanNandi/disvault/app"
	"github.com/AnkanNandi/disvault/db"
	"github.com/spf13/cobra"
)

// required flags
var (
	searchText string
	fileID     int
	gID        int
)

// listCmd represents the list command in the file
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the uploaded files",
	Long: `List command shows the uploaded files in a minimal table format.
You may need to use list command to check the ID or group of a file
for download, delete and updating it.`,
	Run: runListCmd,
}

func init() {
	listCmd.Flags().StringVarP(&searchText, "search", "s", "", "Search by file name, put the keywords")
	listCmd.Flags().IntVarP(&fileID, "id", "i", 0, "the exact file ID you may wanna search")
	listCmd.Flags().IntVarP(&gID, "group", "g", 0, "Search by files by group ID")

	listCmd.MarkFlagsMutuallyExclusive("id", "search")

	rootCmd.AddCommand(listCmd)
}

// Main function responsible for the listing of files
func runListCmd(cmd *cobra.Command, args []string) {
	db.InitDatabase()
	app.Init()
	if cmd.Flags().Changed("group") {
		ValidateGroupID(gID)
	}
	filesList, err := fetchFiles(searchText, fileID, gID)
	if err != nil {
		log.Fatalf("error while fetching files: %v", err)
	}
	listAllFiles(filesList)
}

// formatBytes converts bytes to a human-readable string with appropriate units (B, KB, MB, GB).
func formatBytes(bytes int) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

type listFile struct {
	id        int
	name      string
	size      int // Size in bytes
	parts     int
	groupName string
}

// listAllFiles displays a formatted list of files in a tabular format using the provided slice of listFile structs.
// If the slice is empty, it prints a message indicating that no files matched the criteria.
//
// This function is intended to run when no specific flags (such as search, id, or group) are provided by the user,
// displaying all available files up to the query limit.
//
// Parameters:
//   - files: A slice of listFile structs containing details about each file, including file ID, name, size, total parts, and group name.
//
// Behavior:
//   - Prints the files in a tabular format with headers for File ID, File Name, File Size, Total Parts, and File Group.
//   - If the files slice is empty, prints "No files matched your criteria." to inform the user that there were no results.
//
// Example:
//
//	listAllFiles(files)
//	This will output the list of files in a nicely formatted table or a message if the list is empty.
func listAllFiles(files []listFile) {
	// Check if the files slice is empty
	if len(files) == 0 {
		fmt.Println("No files matched your criteria.")
		return
	}

	// Create a new tabwriter for formatted output
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	// Print the header
	fmt.Fprintln(writer, "FILE ID\tFILE NAME\tFILE SIZE\tTOTAL PARTS\tFILE GROUP")

	// Print the data rows
	for _, file := range files {
		fmt.Fprintf(writer, "%d\t%s\t%s\t%d\t%s\n", file.id, file.name, formatBytes(file.size), file.parts, file.groupName)
	}

	// Flush the writer to ensure the data is written to the output
	writer.Flush()
}

// fetchFiles retrieves a list of files from the database based on the specified search criteria.
// It supports filtering by file name (using a search pattern), file ID, and group ID.
// The results are limited to a maximum of 50 files.
//
// Parameters:
//   - search: A string pattern to match against file names using SQL LIKE syntax. If empty, no name filtering is applied.
//   - id: An integer representing the file ID to search for. If 0, this filter is ignored.
//   - group: An integer representing the group ID to filter files by. If 0, this filter is ignored.
//
// Returns:
//   - A slice of listFile structs containing the matched files' details, including file ID, name, size, total parts, and group name.
//   - An error if there was an issue executing the query or scanning the results.
//
// Example:
//
//	files, err := fetchFiles("report", 0, 2)
//	This call fetches files whose names contain "report" and belong to group ID 2, limited to 50 entries.
//
// Potential Errors:
//   - Returns an error if the database query fails or if there is an issue scanning the rows.
//   - Specific errors might include database connection issues, syntax errors in the SQL query, or type mismatches during row scanning.
func fetchFiles(search string, id int, group int) ([]listFile, error) {
	// Base query to select files
	query := `
		SELECT f.id, f.name, f.size, f.total_parts, g.group_name
		FROM files f
		LEFT JOIN groups g ON f.group_id = g.group_id
		WHERE 1=1
	`
	// Parameters slice for query arguments
	var params []interface{}

	// Conditional query building based on provided flags
	if search != "" {
		query += " AND f.name LIKE ?"
		params = append(params, "%"+search+"%")
	}
	if id != 0 {
		query += " AND f.id = ?"
		params = append(params, id)
	}
	if group != 0 {
		query += " AND g.group_id = ?"
		params = append(params, group)
	}

	// Limit the results to 50 entries
	query += " LIMIT 50"

	// Execute the query with parameters
	rows, err := db.DB.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error querying files: %w", err)
	}
	defer rows.Close()

	// Slice to hold the results
	var files []listFile

	// Iterate over the rows
	for rows.Next() {
		var file listFile
		err := rows.Scan(&file.id, &file.name, &file.size, &file.parts, &file.groupName)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		files = append(files, file)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return files, nil
}
