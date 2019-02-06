package circleci

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/lonelyelk/what-build/aws"
)

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

// FindBuild looks for a build in CircleCI
func FindBuild(projName string, buildName string) (ciBuild *CIBuildResponse, err error) {
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
