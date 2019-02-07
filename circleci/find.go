package circleci

import (
	"encoding/json"
	"fmt"
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

func getBuilds(projectConfig *aws.Project, offset int) (builds []CIBuildResponse, err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := buildRequest(projectConfig.URL, projectConfig.Token, aws.RemoteConfig.Settings.PerPage, offset)
	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	builds = make([]CIBuildResponse, 0)
	err = json.NewDecoder(res.Body).Decode(&builds)
	return
}

func findByBuildParameters(builds *[]CIBuildResponse, params map[string]interface{}) *CIBuildResponse {
	for _, cib := range *builds {
		if cib.BuildParameters == nil {
			continue
		}
		match := true
		for key, value := range params {
			if cib.BuildParameters[key] != value {
				match = false
				break
			}
		}
		if match {
			return &cib
		}
	}
	return nil
}

// FindBuild looks for a build in CircleCI
func FindBuild(projCfg *aws.Project, buildCfg *aws.Build) (*CIBuildResponse, error) {
	for offset := 0; offset < aws.RemoteConfig.Settings.MaxOffset; offset = offset + aws.RemoteConfig.Settings.PerPage {
		ciBuilds, err := getBuilds(projCfg, offset)
		if err != nil {
			return nil, err
		}
		if cib := findByBuildParameters(&ciBuilds, buildCfg.BuildParameters); cib != nil {
			return cib, nil
		}
		if len(ciBuilds) < aws.RemoteConfig.Settings.PerPage {
			break
		}
	}
	return nil, fmt.Errorf("Build %s not found", buildCfg.Name)
}
