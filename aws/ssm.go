package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/viper"
)

// Settings is a settings object for the crawler
type Settings struct {
	PerPage   int `json:"per_page"`
	MaxOffset int `json:"max_offset"`
}

// BuildParameters is a type of generic parameters map
type BuildParameters map[string]interface{}

// Project contains info to fetch builds from CircleCI
type Project struct {
	Name                 string `json:"name"`
	CircleCIURL          string `json:"circleci_url"`
	CircleCIToken        string `json:"circleci_token"`
	CircleCITokenSSMName string `json:"circleci_token_ssm_name"`
}

// Build contains search conditions and identification
type Build struct {
	Name                  string          `json:"name"`
	SearchBuildParameters BuildParameters `json:"search_build_parameters"`
	RunBuildParameters    BuildParameters `json:"run_build_parameters"`
}

// Config contains projects and builds along with settings for the crawler
type Config struct {
	Settings    `json:"settings"`
	Projects    []Project `json:"projects"`
	Builds      []Build   `json:"builds"`
	IAMUserName string
}

var remoteConfig *Config

// GetRemoteConfig returns cached config from SSM or fetches one
func GetRemoteConfig() *Config {
	if remoteConfig != nil {
		return remoteConfig
	}
	return readConfig()
}

// GetSSMParameter fetches a parameter by name from configured region
func GetSSMParameter(name string) (string, error) {
	region := viper.GetString("aws_region")
	session := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	paramIn := ssm.GetParameterInput{Name: aws.String(name), WithDecryption: aws.Bool(true)}
	ssmc := ssm.New(session)
	paramOut, err := ssmc.GetParameter(&paramIn)
	if err != nil {
		return "", err
	}
	return *paramOut.Parameter.Value, nil
}

// GetIAMUserName returns current IAM user name
func GetIAMUserName() (string, error) {
	region := viper.GetString("aws_region")
	session := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	iamc := iam.New(session)
	userOut, err := iamc.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}
	return *userOut.User.UserName, nil
}

func readConfig() *Config {
	cfg, err := GetSSMParameter(viper.GetString("aws_ssm_configuration"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	remoteConfig = &Config{}
	json.NewDecoder(strings.NewReader(cfg)).Decode(remoteConfig)
	if len(remoteConfig.Builds) == 0 || len(remoteConfig.Builds) == 0 {
		fmt.Println("Error reading SSM config")
		os.Exit(1)
	}
	name, err := GetIAMUserName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	remoteConfig.IAMUserName = name
	return remoteConfig
}
