package what

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"

	"github.com/lonelyelk/what-build/aws"
	"github.com/lonelyelk/what-build/circleci"
	"github.com/lonelyelk/what-build/github"
)

func findOrPromptProject(p string, ps *[]aws.Project) (projCfg *aws.Project) {
	projCfg = aws.FindProject(p, ps)
	if projCfg == nil {
		names := make([]string, len(*ps))
		for i, pc := range *ps {
			names[i] = pc.Name
		}
		prompt := promptui.Select{
			Label: "Select a project",
			Items: names,
		}
		_, name, err := prompt.Run()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return aws.FindProject(name, ps)
	}
	return
}

func findOrPromptBuild(b string, bs *[]aws.Build) (buildCfg *aws.Build) {
	buildCfg = aws.FindBuild(b, bs)
	if buildCfg == nil {
		names := make([]string, len(*bs))
		for i, bc := range *bs {
			names[i] = bc.Name
		}
		prompt := promptui.Select{
			Label: "Select a build",
			Items: names,
		}
		_, name, err := prompt.Run()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return aws.FindBuild(name, bs)
	}
	return
}

func promptOptional(pc *aws.Project) aws.BuildParameters {
	if len(pc.OptionalBuildParameters) == 0 {
		return aws.BuildParameters{}
	}
	options := []string{"default"}
	for opt := range pc.OptionalBuildParameters {
		if opt != "default" {
			options = append(options, opt)
		}
	}
	prompt := promptui.Select{
		Label: "Select build options",
		Items: options,
	}
	_, option, err := prompt.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pc.OptionalBuildParameters[option]
}

// Run triggers CircleCI build of given project with run params of given build config
func Run(project string, build string) {
	config := aws.GetRemoteConfig()
	projCfg := findOrPromptProject(project, &config.Projects)
	buildCfg := findOrPromptBuild(build, &config.Builds)
	branch, err := github.ListAndPromptBranch(projCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	params := promptOptional(projCfg)

	ciBuild, err := circleci.RunBuild(projCfg, buildCfg, branch, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	PrintBuild(buildCfg.Name, ciBuild)
}
