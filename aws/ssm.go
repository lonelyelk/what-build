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

// StringIndent constructs a string to output BuildParameters with and indent for formatting
func (params BuildParameters) StringIndent(indent string) string {
	if params == nil || len(params) == 0 {
		return fmt.Sprintf("%s-", indent)
	}
	arr := make([]string, len(params))
	count := 0
	for key, value := range params {
		arr[count] = fmt.Sprintf("%s%s=%v", indent, key, value)
		count = count + 1
	}
	return strings.Join(arr, "\n")
}

// OptionalBuildParameters is a map of BuildParameters
type OptionalBuildParameters map[string]BuildParameters

// DefaultBuildParametersName is a key for default build parameters (always selected first)
const DefaultBuildParametersName = "default"

// StringIndent constructs a string to output OptionalBuildParameters with and indent for formatting
func (opts OptionalBuildParameters) StringIndent(indent string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s%s:\n%s", indent, DefaultBuildParametersName, opts[DefaultBuildParametersName].StringIndent(fmt.Sprintf("%s  - ", indent)))
	for key, value := range opts {
		if key == DefaultBuildParametersName {
			continue
		}
		fmt.Fprintf(&b, "\n%s%s:\n%s", indent, key, value.StringIndent(fmt.Sprintf("%s  - ", indent)))
	}
	return b.String()
}

// Project contains info to fetch builds from CircleCI
type Project struct {
	Name                    string                  `json:"name"`
	CircleCIURL             string                  `json:"circleci_url"`
	CircleCIToken           string                  `json:"circleci_token"`
	CircleCITokenSSMName    string                  `json:"circleci_token_ssm_name"`
	GitHubURL               string                  `json:"github_url"`
	OptionalBuildParameters OptionalBuildParameters `json:"optional_build_parameters"`
}

// Build contains search conditions and identification
type Build struct {
	Name                  string          `json:"name"`
	SearchBuildParameters BuildParameters `json:"search_build_parameters"`
	RunBuildParameters    BuildParameters `json:"run_build_parameters"`
}

// Config contains projects and builds along with settings for the crawler
type Config struct {
	Settings `json:"settings"`
	Projects []Project `json:"projects"`
	Builds   []Build   `json:"builds"`
}

var remoteConfig *Config
var iamUserName string

// GetRemoteConfig returns cached config from SSM or fetches one
func GetRemoteConfig() *Config {
	if remoteConfig != nil {
		return remoteConfig
	}
	return fetchConfig()
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
func GetIAMUserName() string {
	if iamUserName != "" {
		return iamUserName
	}
	return fetchIAMUserName()
}

// FetchTokenIfMissing gets project token from SSM if SSM name is present but token is not, updates project
func FetchTokenIfMissing(p *Project) {
	if p.CircleCIToken == "" {
		token, err := GetSSMParameter(p.CircleCITokenSSMName)
		if err != nil {
			return
		}
		// Account for post request prepared token like 'token:' as if it was user with no password
		if token[len(token)-1] == ':' {
			token = token[:len(token)-1]
		}
		p.CircleCIToken = token
	}
}

func fetchIAMUserName() string {
	region := viper.GetString("aws_region")
	session := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	iamc := iam.New(session)
	userOut, err := iamc.GetUser(&iam.GetUserInput{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	iamUserName = *userOut.User.UserName
	return iamUserName
}

func fetchConfig() *Config {
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
	return remoteConfig
}

// SatisfiedBy checks if provided map includes all the same values under the same keys as the source
func (bp BuildParameters) SatisfiedBy(other BuildParameters) bool {
	if other == nil {
		return false
	}
	for key, value := range bp {
		if other[key] != value {
			return false
		}
	}
	return true
}
