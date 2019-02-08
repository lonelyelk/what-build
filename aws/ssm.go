package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/viper"
)

// Settings is a settings object for the crawler
type Settings struct {
	PerPage   int `json:"per_page"`
	MaxOffset int `json:"max_offset"`
}

// Project contains info to fetch builds from CircleCI
type Project struct {
	Name          string `json:"name"`
	CircleCIURL   string `json:"circleci_url"`
	CircleCIToken string `json:"circleci_token"`
}

// Build contains search conditions and identification
type Build struct {
	Name                  string                 `json:"name"`
	SearchBuildParameters map[string]interface{} `json:"search_build_parameters"`
}

// Config contains projects and builds along with settings for the crawler
type Config struct {
	Settings Settings  `json:"settings"`
	Projects []Project `json:"projects"`
	Builds   []Build   `json:"builds"`
}

// RemoteConfig of the tool
var RemoteConfig = &Config{}

// ReadConfig needs to be run to have Configuration
func ReadConfig() {
	region := viper.GetString("aws_region")
	ssmPath := viper.GetString("aws_ssm_configuration")
	session := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	// iamc := iam.New(session)
	// userOut, err := iamc.GetUser(&iam.GetUserInput{})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(userOut)
	paramIn := ssm.GetParameterInput{Name: aws.String(ssmPath)}
	ssmc := ssm.New(session)
	paramOut, err := ssmc.GetParameter(&paramIn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	json.NewDecoder(strings.NewReader(*paramOut.Parameter.Value)).Decode(RemoteConfig)
	if len(RemoteConfig.Builds) == 0 || len(RemoteConfig.Builds) == 0 {
		fmt.Println("Error reading SSM config")
		os.Exit(1)
	}
}
