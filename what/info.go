package what

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lonelyelk/what-build/aws"
)

// Info outputs available options for projects and builds
func Info() {
	config := aws.GetRemoteConfig()
	sprintBuild := color.New(color.FgCyan, color.Bold).SprintFunc()
	sprintProject := color.New(color.FgWhite, color.Bold).SprintFunc()
	fmt.Printf("\nAvailable %d projects:\n", len(config.Projects))
	for _, project := range config.Projects {
		fmt.Printf("  - %s\n", sprintProject(project.Name))
		if len(project.OptionalBuildParameters) > 0 {
			fmt.Println(project.OptionalBuildParameters.StringIndent("      "))
		}
	}
	fmt.Printf("\nAvailable %d builds:\n", len(config.Builds))
	for _, build := range config.Builds {
		fmt.Printf("  - %s\n", sprintBuild(build.Name))
	}
	fmt.Println()
}
