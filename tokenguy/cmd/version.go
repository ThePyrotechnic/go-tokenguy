package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version    string
	CommitHash string
	BuildTime  string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
		fmt.Println(CommitHash)
		fmt.Printf("Build time: %s", BuildTime)
	},
}
