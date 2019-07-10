package what

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/codebuild"

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

// PrintBuild prints Circle CI build in beautiful colors
func PrintBuild(buildName string, ciBuild *circleci.CIBuildResponse, codeBuild *codebuild.Build, codeBuildExpected bool) {
	sprintBuild := color.New(color.FgCyan, color.Bold).SprintFunc()
	ifColor := color.FgYellow
	if ciBuild.Status == "success" {
		ifColor = color.FgGreen
	}
	if ciBuild.Status == "failed" {
		ifColor = color.FgRed
	}
	var b strings.Builder
	fmt.Fprintf(&b, "  CircleCI Build %s status %s", sprintBuild(buildName), color.New(ifColor).Sprint(ciBuild.Status))
	if ciBuild.StopTime != "" {
		fmt.Fprintf(&b, fmt.Sprintf(" at %s", timeString(ciBuild.StopTime)))
	}
	if user, ok := ciBuild.BuildParameters["IAM_USER"]; ok {
		fmt.Fprintf(&b, fmt.Sprintf(" by %s", user.(string)))
	}
	fmt.Println(b.String())
	ifColor = color.FgYellow
	if ciBuild.Branch == "master" || ciBuild.Branch == "staging" || ciBuild.Branch == "develop" {
		ifColor = color.FgGreen
	}
	fmt.Printf("    - Branch: %s\n", color.New(ifColor).Sprint(ciBuild.Branch))
	fmt.Printf("    - Commit: %s\n", color.New(color.FgMagenta).Sprint(ciBuild.Subject))
	fmt.Printf("    - Revision: %s\n", color.New(color.FgMagenta).Sprint(ciBuild.VcsRevision))

	if !codeBuildExpected {
		fmt.Println()
		return
	}

	if codeBuild == nil {
		fmt.Printf("  CodeBuild not found\n\n")
	} else {
		ifColor = color.FgYellow
		if *codeBuild.BuildStatus == "SUCCEEDED" {
			ifColor = color.FgGreen
		}
		if *codeBuild.BuildStatus == "FAILED" {
			ifColor = color.FgRed
		}
		fmt.Printf("  CodeBuild status %s\n\n", color.New(ifColor).Sprint(*codeBuild.BuildStatus))
	}
}

// Find looks for CircleCI builds of given projects and prints their info
func Find(projects []string, builds []string) {
	config := aws.GetRemoteConfig()
	projCfgs := aws.FindProjects(projects, &config.Projects)
	buildCfgs := aws.FindBuilds(builds, &config.Builds)

	for _, projCfg := range projCfgs {
		fmt.Printf("\nProject %s:\n", color.New(color.FgWhite, color.Bold).Sprint(projCfg.Name))
		for _, buildCfg := range buildCfgs {
			ciBuild, err := circleci.FindBuild(&projCfg, &buildCfg)
			if err != nil {
				fmt.Println(err)
				break
			}

			codeBuild, err := aws.FindCodeBuild(&projCfg, &buildCfg, ciBuild.VcsRevision)
			if err != nil {
				fmt.Println(err)
				break
			}
			PrintBuild(buildCfg.Name, ciBuild, codeBuild, projCfg.CodeBuildName != "")
		}
	}
}
