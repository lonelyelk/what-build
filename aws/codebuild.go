package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/spf13/viper"
)

// FindCodeBuild returns a corresponding codebuild.Build pointer if found
func FindCodeBuild(projConf *Project, buildConf *Build, sha string) (*codebuild.Build, error) {
	if projConf.CodeBuildName == "" {
		return nil, nil
	}
	region := viper.GetString("aws_region")
	session := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	cbc := codebuild.New(session)
	lbIn := &codebuild.ListBuildsForProjectInput{ProjectName: aws.String(projConf.CodeBuildName)}
	lbOut, err := cbc.ListBuildsForProject(lbIn)
	if err != nil {
		return nil, err
	}
	bgbIn := &codebuild.BatchGetBuildsInput{Ids: lbOut.Ids}
	bgbOut, err := cbc.BatchGetBuilds(bgbIn)
	if err != nil {
		return nil, err
	}
	for _, build := range bgbOut.Builds {
		if buildConf.SearchCodeBuildParameters.SatisfiedByCodeBuildEnv(build.Environment.EnvironmentVariables) {
			if *build.SourceVersion == sha {
				return build, nil
			}
			return nil, nil
		}
	}
	for lbOut.NextToken != nil {
		lbIn.NextToken = lbOut.NextToken
		lbOut, err = cbc.ListBuildsForProject(lbIn)
		if err != nil {
			return nil, err
		}
		bgbIn.Ids = lbOut.Ids
		bgbOut, err = cbc.BatchGetBuilds(bgbIn)
		if err != nil {
			return nil, err
		}
		for _, build := range bgbOut.Builds {
			if buildConf.SearchCodeBuildParameters.SatisfiedByCodeBuildEnv(build.Environment.EnvironmentVariables) {
				if *build.SourceVersion == sha {
					return build, nil
				}
				return nil, nil
			}
		}
	}
	return nil, nil
}
