package circleci

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func errStatus(url *url.URL) error {
	return fmt.Errorf("circleci: project url '%s://%s%s' doesn't exist", url.Scheme, url.Host, url.Path)
}
func errBuildNotFound(buildName string, projectName string) error {
	return fmt.Errorf("circleci: build '%s' not found fpr project '%s'", buildName, projectName)
}

// FetchBuildsRequest constructs and returns CircleCI API based request for builds
func FetchBuildsRequest(url string, token string, limit int, offset int) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("circle-token", token)
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lonelyelk-what-build-bot")
	req.URL.RawQuery = q.Encode()
	return
}

// FetchBuildsDo makes a request to a URL from project config to fetch Circle CI builds with limit and offset
func FetchBuildsDo(projectConfig *aws.Project, limit int, offset int) (builds []CIBuildResponse, err error) {
	req, err := FetchBuildsRequest(projectConfig.CircleCIURL, projectConfig.CircleCIToken, limit, offset)
	if err != nil {
		return
	}
	client := http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			var r = req
			if len(via) > 0 {
				r = via[len(via)-1]
			}
			return errStatus(r.URL)
		},
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return builds, errStatus(req.URL)
	}
	builds = make([]CIBuildResponse, 0)
	err = json.NewDecoder(res.Body).Decode(&builds)
	return
}

// FindByBuildParameters returns the first build where searched parameters have the same values
func FindByBuildParameters(builds *[]CIBuildResponse, params map[string]interface{}) *CIBuildResponse {
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
	config := aws.GetRemoteConfig()
	for offset := 0; offset < config.Settings.MaxOffset; offset = offset + config.Settings.PerPage {
		ciBuilds, err := FetchBuildsDo(projCfg, config.Settings.PerPage, offset)
		if err != nil {
			return nil, err
		}
		if cib := FindByBuildParameters(&ciBuilds, buildCfg.SearchBuildParameters); cib != nil {
			return cib, nil
		}
		if len(ciBuilds) < config.Settings.PerPage {
			break
		}
	}
	return nil, errBuildNotFound(buildCfg.Name, projCfg.Name)
}
