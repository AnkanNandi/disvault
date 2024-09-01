package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/AnkanNandi/disvault/app"
	"github.com/AnkanNandi/disvault/db"
	"github.com/spf13/cobra"
)

// Flags for the group command
var (
	group           string
	parentGroup     string
	deleteGroupName string
	listGroups      bool
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups within DisVault",
	Long:  `Group command allows you to create, delete, and manage groups within DisVault.`,
	Run:   runGroupCmd,
}

func init() {
	groupCmd.Flags().StringVarP(&group, "name", "n", "", "Create a new group with the specified name")
	groupCmd.Flags().StringVarP(&parentGroup, "parent", "p", "", "Specify a parent group by name or ID when creating a new group")
	groupCmd.Flags().StringVarP(&deleteGroupName, "delete", "d", "", "Delete a group using its group name")
	groupCmd.Flags().BoolVarP(&listGroups, "list", "l", false, "List all available groups")
	groupCmd.MarkFlagsMutuallyExclusive("name", "delete", "list")

	rootCmd.AddCommand(groupCmd)
}

func runGroupCmd(cmd *cobra.Command, args []string) {
	db.InitDatabase()
	app.Init()

	if cmd.Flags().Changed("parent") && !cmd.Flags().Changed("name") {
		fmt.Println("Error: The -p (parent group) flag can only be used with the -n (name) flag.\nYou may haven't provided a name.")
		return
	}

	switch {
	case listGroups:
		listAllGroups()
	case group != "":
		createGroup()
	case deleteGroupName != "":
		deleteGroup(deleteGroupName)
	default:
		fmt.Println("Please provide a valid flag. Use -n to create a group or -d to delete a group.")
	}
}

func listAllGroups() {
	// Fetch group details along with their parent group names using a LEFT JOIN
	rows, err := db.DB.Query(`
		SELECT g.group_id, g.group_name, pg.group_name AS parent_group_name
		FROM groups g
		LEFT JOIN groups pg ON g.parent_group_id = pg.group_id
	`)
	if err != nil {
		log.Fatalf("Error fetching groups: %v", err)
	}
	defer rows.Close()

	// Create a new tabwriter for formatted output
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	// Print the header
	fmt.Fprintln(writer, "GROUP ID\tGROUP NAME\tPARENT GROUP NAME")

	// Iterate over the rows and print each group
	for rows.Next() {
		var groupID int
		var groupName string
		var parentGroupName sql.NullString

		// Scan the row into variables
		if err := rows.Scan(&groupID, &groupName, &parentGroupName); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		// Handle NULL values for parent_group_name
		parentName := "NONE/ROOT"
		if parentGroupName.Valid {
			parentName = parentGroupName.String
		}

		// Print the row
		fmt.Fprintf(writer, "%d\t%s\t%s\n", groupID, groupName, parentName)
	}

	// Flush the writer to ensure the data is written to the output
	writer.Flush()

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
	}
}

func createGroup() {
	var parentID interface{} = nil

	if parentGroup != "" {
		id, err := fetchGroupIDByName(parentGroup)
		if err != nil {
			fmt.Println(err)
			return
		}
		parentID = id
	}

	_, err := db.DB.Exec("INSERT INTO groups (group_name, parent_group_id) VALUES (?, ?)", group, parentID)
	if err != nil {
		if isUniqueConstraintError(err) {
			fmt.Printf("Error: A group with the name '%s' already exists. Please choose a different name.\n", group)
		} else {
			log.Fatalf("Error creating group: %v", err)
		}
		return
	}

	fmt.Printf("Group '%s' created successfully.\n", group)
}

func fetchGroupIDByName(name string) (int, error) {
	var groupID int
	err := db.DB.QueryRow("SELECT group_id FROM groups WHERE group_name = ?", name).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("error: Parent group '%s' not found", name)
		}
		log.Fatalf("Error fetching group ID by name: %v", err)
	}
	return groupID, nil
}

func deleteGroup(groupName string) {
	groupID, err := fetchGroupIDByName(groupName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Reassign files to the 'uncategorized' group (group_id = 1)
	_, err = db.DB.Exec("UPDATE files SET group_id = 1 WHERE group_id = ?", groupID)
	if err != nil {
		log.Fatalf("Error reassigning files to 'uncategorized': %v", err)
	}

	// Recursively delete all child groups
	deleteChildGroups(groupID)

	// Delete the parent group from the database
	_, err = db.DB.Exec("DELETE FROM groups WHERE group_id = ?", groupID)
	if err != nil {
		log.Fatalf("Error deleting group: %v", err)
	}

	fmt.Printf("Group '%s' and its child groups deleted successfully, and files reassigned to 'uncategorized'.\n", groupName)
}

func deleteChildGroups(parentID int) {
	rows, err := db.DB.Query("SELECT group_id FROM groups WHERE parent_group_id = ?", parentID)
	if err != nil {
		log.Fatalf("Error fetching child groups: %v", err)
	}
	defer rows.Close()

	var childGroupIDs []int
	for rows.Next() {
		var childID int
		if err := rows.Scan(&childID); err != nil {
			log.Fatalf("Error scanning child group ID: %v", err)
		}
		childGroupIDs = append(childGroupIDs, childID)
	}

	for _, childID := range childGroupIDs {
		// Recursively delete children of the current child group
		deleteChildGroups(childID)

		// Reassign files to the 'uncategorized' group (group_id = 1)
		_, err := db.DB.Exec("UPDATE files SET group_id = 1 WHERE group_id = ?", childID)
		if err != nil {
			log.Fatalf("Error reassigning files to 'uncategorized' for group %d: %v", childID, err)
		}

		// Delete the child group
		_, err = db.DB.Exec("DELETE FROM groups WHERE group_id = ?", childID)
		if err != nil {
			log.Fatalf("Error deleting child group %d: %v", childID, err)
		}

		fmt.Printf("Child group with ID '%d' deleted successfully.\n", childID)
	}
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	sqlErr, ok := err.(interface{ Error() string })
	return ok && strings.Contains(sqlErr.Error(), "UNIQUE constraint failed")
}
