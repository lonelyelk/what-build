package cmd

import (
	"github.com/lonelyelk/what-build/what"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "list available projects and builds",
	Long: `
Read SSM config and list available projects and builds by name, that can be used with flags
`,
	Run: func(cmd *cobra.Command, args []string) {
		what.Info()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
