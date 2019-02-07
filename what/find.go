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
	projCfgs := aws.FindProjects(projects)
	buildCfgs := aws.FindBuilds(builds)

	for _, projCfg := range projCfgs {
		fmt.Printf("\nProject %s:\n", color.New(color.FgWhite, color.Bold).Sprint(projCfg.Name))
		for _, buildCfg := range buildCfgs {
			ciBuild, err := circleci.FindBuild(projCfg, buildCfg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			time, err := time.Parse(time.RFC3339, ciBuild.StopTime)
			sprintBuild := color.New(color.FgCyan, color.Bold).SprintFunc()
			if err != nil {
				fmt.Printf("  Deployed to %s at %v\n", sprintBuild(buildCfg.Name), ciBuild.StopTime)
			} else {
				fmt.Printf("  Deployed to %s at %v\n", sprintBuild(buildCfg.Name), time)
			}
			branchColor := color.FgYellow
			greenBranches := [3]string{"master", "staging", "develop"}
			for _, name := range greenBranches {
				if name == ciBuild.Branch {
					branchColor = color.FgGreen
				}
			}
			fmt.Printf("    - Branch: %s\n", color.New(branchColor).Sprint(ciBuild.Branch))
			fmt.Printf("    - Commit: %s\n", color.New(color.FgMagenta).Sprint(ciBuild.Subject))
			fmt.Printf("    - Revision: %s\n\n", color.New(color.FgMagenta).Sprint(ciBuild.VcsRevision))
		}
	}
}
