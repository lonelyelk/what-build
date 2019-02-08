package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is current tool semantic version
var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of what-build tool",
	Long: `
Current version of a tool
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
