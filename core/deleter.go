package core

import (
	"context"
	"fmt"

	"github.com/AnkanNandi/disvault/app"
	"github.com/AnkanNandi/disvault/db"
	"github.com/bwmarrin/discordgo"
)

/*
The function deletes all the file parts of a file then finally deletes the file itself

It deletes the file first then it's registration in the parts table.

	'TODO: Add a way to check for sudden stops in delete i.e. program crash in middle of function'
*/
func DeleteFileParts(fileID int) error {
	ctx := context.Background()
	partIDs, err := db.GetPartIDs(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to retrieve part IDs: %w", err)
	}

	discord, err := discordgo.New("Bot " + app.Config.BotToken)
	if err != nil {
		return fmt.Errorf("invalid bot parameters: %w", err)
	}
	defer discord.Close()

	for _, partID := range partIDs {
		// Delete the message (file) from Discord
		err := discord.ChannelMessageDelete(app.Config.ChannelID, partID)
		if err != nil {
			return fmt.Errorf("failed to delete message %s from Discord: %w", partID, err)
		}
		fmt.Printf("Deleted part %s from Discord\n", partID)

		// Delete the part entry from the database
		_, err = db.DB.ExecContext(ctx, "DELETE FROM parts WHERE part_id = ?", partID)
		if err != nil {
			return fmt.Errorf("failed to delete part %s from database: %w", partID, err)
		}
		fmt.Printf("Deleted part %s from database\n", partID)
	}

	// also delete the file entry from the 'files' table
	_, err = db.DB.ExecContext(ctx, "DELETE FROM files WHERE id = ?", fileID)
	if err != nil {
		return fmt.Errorf("failed to delete file ID %d from database: %w", fileID, err)
	}
	fmt.Printf("Deleted file ID %d from database\n", fileID)

	return nil
}
