package what

import (
	"fmt"

	"github.com/lonelyelk/what-build/aws"
	"github.com/lonelyelk/what-build/circleci"
)

// Run triggers CircleCI build of given project with run params of given build config
func Run(project string, build string) {
	config := aws.GetRemoteConfig()
	projCfg := aws.FindProject(project, &config.Projects)
	buildCfg := aws.FindBuild(build, &config.Builds)

	ciBuild, err := circleci.RunBuild(projCfg, buildCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	PrintBuild(buildCfg.Name, ciBuild)
}
