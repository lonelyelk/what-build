package what

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/lonelyelk/what-build/aws"
	"github.com/lonelyelk/what-build/circleci"
)

func timeString(str string) string {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return str
	}
	return t.In(time.Local).Format("15:04 02.01.2006")
}

func printBuild(buildName string, ciBuild *circleci.CIBuildResponse) {
	sprintBuild := color.New(color.FgCyan, color.Bold).SprintFunc()
	ifColor := color.FgYellow
	if ciBuild.Status == "success" {
		ifColor = color.FgGreen
	}
	if ciBuild.Status == "failed" {
		ifColor = color.FgRed
	}
	fmt.Printf("  Build %s status %s", sprintBuild(buildName), color.New(ifColor).Sprint(ciBuild.Status))
	if ciBuild.StopTime != "" {
		fmt.Printf(" at %s\n", timeString(ciBuild.StopTime))
	}
	ifColor = color.FgYellow
	if ciBuild.Branch == "master" || ciBuild.Branch == "staging" || ciBuild.Branch == "develop" {
		ifColor = color.FgGreen
	}
	fmt.Printf("    - Branch: %s\n", color.New(ifColor).Sprint(ciBuild.Branch))
	fmt.Printf("    - Commit: %s\n", color.New(color.FgMagenta).Sprint(ciBuild.Subject))
	fmt.Printf("    - Revision: %s\n\n", color.New(color.FgMagenta).Sprint(ciBuild.VcsRevision))
}

// Find looks for CircleCI builds of given projects and prints their info
func Find(projects []string, builds []string) {
	projCfgs := aws.FindProjects(projects)
	buildCfgs := aws.FindBuilds(builds)

	for _, projCfg := range projCfgs {
		fmt.Printf("\nProject %s:\n", color.New(color.FgWhite, color.Bold).Sprint(projCfg.Name))
		for _, buildCfg := range buildCfgs {
			ciBuild, err := circleci.FindBuild(&projCfg, &buildCfg)
			if err != nil {
				fmt.Println(err)
				break
			}
			printBuild(buildCfg.Name, ciBuild)
		}
	}
}
