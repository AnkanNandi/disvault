package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the application version, set via ldflags during build.
var Version = "0.0.1-alpha.1" // Default value for development builds

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of disvault",
	Long:  `All software has versions. This is DisVault's version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("DisVault version: %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
