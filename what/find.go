package what

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/lonelyelk/what-build/aws"
	"github.com/lonelyelk/what-build/circleci"
)

// Find looks for CircleCI builds of given projects and prints their info
func Find(projects []string, builds []string) {
	config := aws.RemoteConfig
	if len(config.Builds) == 0 || len(config.Builds) == 0 {
		fmt.Println("Error reading SSM config")
		os.Exit(1)
	}

	if len(projects) == 0 {
		projects = make([]string, len(config.Projects))
		for i, p := range config.Projects {
			projects[i] = p.Name
		}
	}

	if len(builds) == 0 {
		builds = make([]string, len(config.Builds))
		for i, b := range config.Builds {
			builds[i] = b.Name
		}
	}

	for _, project := range projects {
		fmt.Printf("\nProject %s:\n", color.New(color.FgWhite, color.Bold).Sprint(project))
		for _, build := range builds {
			ciBuild, err := circleci.FindBuild(project, build)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			time, err := time.Parse(time.RFC3339, ciBuild.StopTime)
			sprintBuild := color.New(color.FgCyan, color.Bold).SprintFunc()
			if err != nil {
				fmt.Printf("  Deployed to %s at %v\n", sprintBuild(build), ciBuild.StopTime)
			} else {
				fmt.Printf("  Deployed to %s at %v\n", sprintBuild(build), time)
			}
			branchColor := color.FgYellow
			greenBranches := [3]string{"master", "staging", "develop"}
			for _, name := range greenBranches {
				if name == ciBuild.Branch {
					branchColor = color.FgGreen
				}
			}
			fmt.Printf("    - Branch: %s\n", color.New(branchColor).Sprint(ciBuild.Branch))
			fmt.Printf("    - Commit: %s\n", color.New(color.FgBlue).Sprint(ciBuild.Subject))
			fmt.Printf("    - Revision: %s\n\n", color.New(color.FgMagenta).Sprint(ciBuild.VcsRevision))
		}
	}
}
