package cmd

import (
	"github.com/lonelyelk/what-build/what"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of what-build tool",
	Long: `
Current version of a tool
`,
	Run: func(cmd *cobra.Command, args []string) {
		what.PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
