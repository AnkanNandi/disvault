package core

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/AnkanNandi/disvault/app"
	"github.com/AnkanNandi/disvault/db"
)

const chunkSize = 25 * 1024 * 1024 // 25 MB

func Upload(inputFile string, groupID int) error {
	ctx := context.Background()

	// Create a temporary directory for file chunks
	tempDir, err := os.MkdirTemp("", "temp")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Open the input file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	// Retrieve file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Calculate file hash
	fileHash, err := FileHash(ctx, file)
	if err != nil {
		return fmt.Errorf("failed to calculate file hash: %w", err)
	}

	// Reset file pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Register the main file in the database
	fileToBeUploaded := db.FilesLocal{
		Name:        fileInfo.Name(),
		Total_parts: FilePartsCalc(fileInfo.Size()),
		Size:        fileInfo.Size(),
		Hash:        fileHash,
		GroupID:     groupID,
	}

	mainFileID, err := db.RegisterFileEntry(ctx, &fileToBeUploaded)
	if err != nil {
		return fmt.Errorf("failed to register file in database: %w", err)
	}

	// Buffer for reading file chunks
	buffer := make([]byte, chunkSize)
	fileParts := db.Parts{FileID: int(mainFileID)}

	// Read and upload chunks
	for i := 0; ; i++ {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading file chunk: %w", err)
		}
		if bytesRead == 0 {
			break
		}

		// Write chunk to temporary file
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", fileInfo.Name(), i))
		if err := os.WriteFile(chunkPath, buffer[:bytesRead], 0644); err != nil {
			return fmt.Errorf("error writing chunk file: %w", err)
		}

		// Upload chunk
		if _, err := app.UploadFile(ctx, chunkPath, &fileParts); err != nil {
			return fmt.Errorf("error uploading chunk: %w", err)
		}

		// Stop if at the end of file
		if bytesRead < chunkSize {
			break
		}
	}

	// Insert file parts into the database
	if err := fileParts.InsertParts(ctx, fileParts); err != nil {
		return fmt.Errorf("error inserting file parts into database: %w", err)
	}

	fmt.Println("File upload completed successfully.")
	return nil
}

// FilePartsCalc calculates the number of parts required to split the file.
func FilePartsCalc(fileSize int64) int {
	return int(math.Ceil(float64(fileSize) / chunkSize))
}

// FileHash calculates the SHA-256 hash of the given file.
func FileHash(ctx context.Context, file *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("error hashing file: %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
