package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	DB   *sql.DB
	once sync.Once
)

// Groups table contains data for group so it can be joined when called by the user
// files table contains all the metadata of the file and how many parts it may have,
// that total_parts is calculated by a function and doesn't care about the discord files.
// Parts table has all the parts that are uploaded by the UploadFiles function
// On the groups table, root group means if a group is related to some other group,
// TODO: add better names
const Tables string = `-- Create the 'groups' table with a self-referencing foreign key
			CREATE TABLE IF NOT EXISTS groups (
   			 group_id INTEGER PRIMARY KEY AUTOINCREMENT,
   			 group_name TEXT UNIQUE NOT NULL,
   			 parent_group_id INTEGER,  -- Self-referencing column
    		FOREIGN KEY (parent_group_id) REFERENCES groups(group_id)
			);

			-- Create a unique index for 'group_name' in 'groups'
			CREATE UNIQUE INDEX IF NOT EXISTS idx_group_name ON groups(group_name);

			-- Insert the default group 'uncategorized' only if it doesn't exist
			INSERT INTO groups (group_name)
			SELECT 'uncategorized'
			WHERE NOT EXISTS (SELECT 1 FROM groups WHERE group_name = 'uncategorized');


			-- Create the 'files' table
			CREATE TABLE IF NOT EXISTS files (
   			 id INTEGER PRIMARY KEY AUTOINCREMENT,
    		 name TEXT NOT NULL,
    		 total_parts INTEGER NOT NULL,
   			 size INTEGER NOT NULL,
   			 hash TEXT NOT NULL,
    		 group_id INTEGER NOT NULL DEFAULT 1,  -- Set default group to 'uncategorized'
    		 FOREIGN KEY (group_id) REFERENCES groups(group_id)
			);

			-- Create indexes for 'files'
			CREATE INDEX IF NOT EXISTS idx_file_search_id ON files(group_id);
			CREATE INDEX IF NOT EXISTS idx_file_search_name ON files(name);

			-- Create the 'parts' table
			CREATE TABLE IF NOT EXISTS parts (
   			 part_id TEXT PRIMARY KEY,
   			 file_id INTEGER NOT NULL,
   			 FOREIGN KEY (file_id) REFERENCES files(id)
			);

			-- Create index for 'parts'
			CREATE INDEX IF NOT EXISTS idx_file_id ON parts(file_id);
`

// InitDatabase initializes the database, creating necessary tables if they don't exist.
func InitDatabase() error {
	var err error

	once.Do(func() {
		// Create the output directory
		dbPath := filepath.Join(".", "data")
		err = os.MkdirAll(dbPath, 0755)
		// Open the database connection
		DB, err = sql.Open("sqlite", filepath.Join(dbPath, "db.sql"))
		if err != nil {
			err = fmt.Errorf("failed to open database: %w", err)
			return
		}

		// Ensure the database is accessible
		if err = DB.Ping(); err != nil {
			err = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		// Create necessary tables
		_, err = DB.ExecContext(
			context.Background(),
			Tables,
		)
		if err != nil {
			err = fmt.Errorf("failed to create tables: %w", err)
			return
		}
	})

	return err
}

// RegisterFileEntry adds a file entry to the files table in the database.
func RegisterFileEntry(ctx context.Context, fileStructure *FilesLocal) (int64, error) {
	result, err := DB.ExecContext(
		ctx,
		"INSERT INTO files (name, total_parts, size, hash, group_id) VALUES (?, ?, ?, ?, ?)",
		fileStructure.Name, fileStructure.Total_parts, fileStructure.Size, fileStructure.Hash, fileStructure.GroupID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to register file: %w", err)
	}

	fileID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last inserted ID: %w", err)
	}

	fmt.Printf("The file ID: %v\n", fileID)
	return fileID, nil
}

// InsertParts Allows for registering the uploaded files in the database
/* (data Parts) is the []string of data of registered files */
func (parts *Parts) InsertParts(ctx context.Context, data Parts) error {
	// Prepare the insert statement
	stmt, err := DB.PrepareContext(ctx, "INSERT INTO parts (part_id, file_id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Iterate over each part and insert into the database
	for _, partID := range data.Parts {
		_, err := stmt.ExecContext(ctx, partID, data.FileID)
		if err != nil {
			log.Printf("failed to insert part: %s for file ID: %d, error: %v", partID, data.FileID, err)
			return err // the function will stop
		}
		fmt.Printf("Successfully inserted part: %s for file ID: %d\n", partID, data.FileID)
	}

	return nil
}

// gets all the fileparts for use by different function,
// it returns them in orders since discord uses a timestamp as id it works fine
func GetPartIDs(ctx context.Context, fileID int) ([]string, error) {
	rows, err := DB.QueryContext(ctx, "SELECT part_id FROM parts WHERE file_id = ? ORDER BY part_id", fileID)
	if err != nil {
		return nil, fmt.Errorf("error querying parts: %w", err)
	}
	defer rows.Close()

	var partIDs []string
	for rows.Next() {
		var partID string
		if err := rows.Scan(&partID); err != nil {
			return nil, fmt.Errorf("error scanning part_id: %w", err)
		}
		partIDs = append(partIDs, partID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return partIDs, nil
}

// FilesLocal represents a local file's metadata.
type FilesLocal struct {
	Name        string // Name of the file
	Total_parts int    // Number of parts the file is divided into
	Size        int64  // Size of the file
	Hash        string // Hash for verifying file integrity, on complete download user may check for hash match
	GroupID     int    // User may assign group to each file for future search commands, multiple files may belong to same group i.e. math books
}

// FilesDB represents a file entry in the database, including its ID.
type FilesDB struct {
	ID int // Auto-incremented ID, PRIMARY KEY
	FilesLocal
}

// Used to register the file parts to the DB
type Parts struct {
	FileID int      // File ID that was registered earlier
	Parts  []string // The parts message IDs uploaded by the app.UploadFile() function
}
