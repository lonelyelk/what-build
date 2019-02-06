package what

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/lonelyelk/what-build/aws"
)

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
func findCIBuild(projName string, buildName string) (ciBuild *CIBuildResponse, err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	config := aws.RemoteConfig
	var buildConfig *aws.Build
	var projectConfig *aws.Project
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

// Find looks for CircleCI builds of given projects and prints their info
func Find(projects []string, builds []string) {
	config := aws.RemoteConfig
	if len(config.Builds) == 0 || len(config.Builds) == 0 {
		fmt.Println("Error reading SSM config")
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
		fmt.Printf("\nProject %s:\n", color.New(color.FgWhite, color.Bold).Sprint(project))
		for _, build := range builds {
			ciBuild, err := findCIBuild(project, build)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			time, err := time.Parse(time.RFC3339, ciBuild.StopTime)
			sprintBuild := color.New(color.FgCyan, color.Bold).SprintFunc()
			if err != nil {
				fmt.Printf("  Deployed to %s at %v\n", sprintBuild(build), ciBuild.StopTime)
			} else {
				fmt.Printf("  Deployed to %s at %v\n", sprintBuild(build), time)
			}
			branchColor := color.FgYellow
			greenBranches := [3]string{"master", "staging", "develop"}
			for _, name := range greenBranches {
				if name == ciBuild.Branch {
					branchColor = color.FgGreen
				}
			}
			fmt.Printf("    - Branch: %s\n", color.New(branchColor).Sprint(ciBuild.Branch))
			fmt.Printf("    - Commit: %s\n", color.New(color.FgBlue).Sprint(ciBuild.Subject))
			fmt.Printf("    - Revision: %s\n\n", color.New(color.FgMagenta).Sprint(ciBuild.VcsRevision))
		}
	}
}
