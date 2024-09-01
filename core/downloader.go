package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AnkanNandi/disvault/app"
	"github.com/AnkanNandi/disvault/db"
)

// Assembles the binary files in one big file
func DownloadAndReassembleFile(fileID int, outputFileName string) error {
	ctx := context.Background()
	partIDs, err := db.GetPartIDs(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to retrieve part IDs: %w", err)
	}

	// Create the output directory
	outputDir := filepath.Join(".", "out")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the output file
	outputFilePath := filepath.Join(outputDir, outputFileName)
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Download and stitch each part together
	for i, partID := range partIDs {
		fmt.Printf("Downloading part %d/%d: %s\n", i+1, len(partIDs), partID)
		partData, err := app.DownloadPart(partID)
		if err != nil {
			return fmt.Errorf("failed to download part %s: %w", partID, err)
		}

		_, err = outFile.Write(partData)
		if err != nil {
			return fmt.Errorf("failed to write part %s to output file: %w", partID, err)
		}
	}

	fmt.Printf("Successfully reassembled file to %s\n", outputFilePath)
	return nil
}
