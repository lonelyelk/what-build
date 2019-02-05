package cmd

import (
	"github.com/lonelyelk/what-build/what"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find [flags]",
	Short: "find a build of a project",
	Long: `
Find a build in CircleCI with Build Parameters configured in SSM config
`,
	Run: func(cmd *cobra.Command, args []string) {
		projects, _ := cmd.Flags().GetStringSlice("project")
		builds, _ := cmd.Flags().GetStringSlice("build")
		what.Find(projects, builds)
	},
}

func init() {
	findCmd.Flags().StringSliceP("project", "p", []string{}, "project or comma-separated project list")
	findCmd.Flags().StringSliceP("build", "b", []string{}, "build or comma-separated build list")
	rootCmd.AddCommand(findCmd)
}
