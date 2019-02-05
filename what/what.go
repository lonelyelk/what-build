package what

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

// Settings is a settings object for the crawler
type Settings struct {
	PerPage   int `json:"per_page"`
	MaxOffset int `json:"max_offset"`
}

// Project contains info to fetch builds from CircleCI
type Project struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

// Build contains search conditions and identification
type Build struct {
	Name            string                 `json:"name"`
	BuildParameters map[string]interface{} `json:"build_parameters"`
}

// Config contains projects and builds along with settings for the crawler
type Config struct {
	Settings Settings  `json:"settings"`
	Projects []Project `json:"projects"`
	Builds   []Build   `json:"builds"`
}

// CIBuildResponse is a JSON extraits for build entity on circleci
type CIBuildResponse struct {
	BuildNum        int                    `json:"build_num"`
	Branch          string                 `json:"branch"`
	VcsRevision     string                 `json:"vcs_revision"`
	Subject         string                 `json:"subject"`
	Why             string                 `json:"why"`
	DontBuild       string                 `json:"dont_build"`
	StopTime        string                 `json:"stop_time"`
	BuildTimeMillis int                    `json:"build_time_millis"`
	Status          string                 `json:"status"`
	BuildParameters map[string]interface{} `json:"build_parameters"`
}

func buildRequest(url string, token string, limit int, offset int) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("circle-token", token)
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	q.Add("filter", "completed")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lonelyelk-what-build-bot")
	req.URL.RawQuery = q.Encode()
	return
}
func findCIBuild(projName string, buildName string, config *Config) (ciBuild *CIBuildResponse, err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	var buildConfig *Build
	var projectConfig *Project
	for _, p := range config.Projects {
		if p.Name == projName {
			projectConfig = &p
			break
		}
	}
	for _, b := range config.Builds {
		if b.Name == buildName {
			buildConfig = &b
			break
		}
	}
	if projectConfig == nil {
		return nil, errors.New("Project config not found")
	}
	if buildConfig == nil {
		return nil, errors.New("Build config not found")
	}
	for offset := 0; offset < config.Settings.MaxOffset; offset = offset + config.Settings.PerPage {
		req, err := buildRequest(projectConfig.URL, projectConfig.Token, config.Settings.PerPage, offset)
		if err != nil {
			return nil, err
		}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		ciBuilds := make([]CIBuildResponse, 0)
		if err = json.NewDecoder(res.Body).Decode(&ciBuilds); err != nil {
			return nil, err
		}
		for _, cib := range ciBuilds {
			if cib.BuildParameters == nil {
				continue
			}
			match := true
			for key, value := range buildConfig.BuildParameters {
				if cib.BuildParameters[key] != value {
					match = false
					break
				}
			}
			if match {
				return &cib, nil
			}
		}
		if len(ciBuilds) < config.Settings.PerPage {
			break
		}
	}
	return nil, errors.New("Build not found")
}

func getConfig() (config *Config, err error) {
	awsRegion := viper.GetString("aws_region")
	awsSsmPath := viper.GetString("aws_ssm_configuration")
	session := session.Must(session.NewSession(&aws.Config{Region: aws.String(awsRegion)}))
	// iamc := iam.New(session)
	// userOut, err := iamc.GetUser(&iam.GetUserInput{})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(userOut)
	paramIn := ssm.GetParameterInput{Name: aws.String(awsSsmPath)}
	ssmc := ssm.New(session)
	paramOut, err := ssmc.GetParameter(&paramIn)
	if err != nil {
		return
	}
	config = &Config{}
	if err = json.NewDecoder(strings.NewReader(*paramOut.Parameter.Value)).Decode(config); err != nil {
		return nil, err
	}
	return
}

// Find looks for CircleCI builds of given projects and prints their info
func Find(projects []string, builds []string) {
	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
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
		fmt.Printf("\nProject %s:\n", project)
		for _, build := range builds {
			ciBuild, err := findCIBuild(project, build, config)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			time, err := time.Parse(time.RFC3339, ciBuild.StopTime)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("  Deployed to %s at %v\n", build, time)
			fmt.Printf("    - Branch: %s\n", ciBuild.Branch)
			fmt.Printf("    - Commit: %s\n", ciBuild.Subject)
			fmt.Printf("    - Revision: %s\n\n", ciBuild.VcsRevision)
		}
	}
}