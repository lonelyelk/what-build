package what

import (
	"fmt"

	"github.com/lonelyelk/what-build/aws"
)

// Info outputs available options for projects and builds
func Info() {
	config := aws.GetRemoteConfig()
	fmt.Printf("\nAvailable %d projects:\n", len(config.Projects))
	for _, project := range config.Projects {
		fmt.Printf("  - %s\n", project.Name)
	}
	fmt.Printf("\nAvailable %d builds:\n", len(config.Builds))
	for _, build := range config.Builds {
		fmt.Printf("  - %s\n", build.Name)
	}
	fmt.Println()
}
