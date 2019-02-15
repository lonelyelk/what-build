package cmd

import (
	"github.com/lonelyelk/what-build/what"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [flags]",
	Short: "run a build of a project",
	Long: `
Run a build in CircleCI with Build Parameters configured in SSM config
`,
	Run: func(cmd *cobra.Command, args []string) {
		projects, _ := cmd.Flags().GetString("project")
		builds, _ := cmd.Flags().GetString("build")
		what.Run(projects, builds)
	},
}

func init() {
	runCmd.Flags().StringP("project", "p", "", "project")
	runCmd.Flags().StringP("build", "b", "", "build")
	rootCmd.AddCommand(runCmd)
}
