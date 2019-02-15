package circleci

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lonelyelk/what-build/api"
	"github.com/lonelyelk/what-build/aws"
)

// CIBuildResponse is a JSON extraits for build entity on circleci
type CIBuildResponse struct {
	BuildNum        int                 `json:"build_num"`
	Branch          string              `json:"branch"`
	VcsRevision     string              `json:"vcs_revision"`
	Subject         string              `json:"subject"`
	Why             string              `json:"why"`
	DontBuild       string              `json:"dont_build"`
	StopTime        string              `json:"stop_time"`
	BuildTimeMillis int                 `json:"build_time_millis"`
	Status          string              `json:"status"`
	BuildParameters aws.BuildParameters `json:"build_parameters"`
}

func errBuildNotFound(buildName string, projectName string) error {
	return fmt.Errorf("circleci: build '%s' not found for project '%s'", buildName, projectName)
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
	req.URL.RawQuery = q.Encode()
	return
}

// FetchBuildsDo makes a request to a URL from project config to fetch Circle CI builds with limit and offset
func FetchBuildsDo(projConfig *aws.Project, limit int, offset int) (builds []CIBuildResponse, err error) {
	aws.FetchTokenIfMissing(projConfig)
	req, err := FetchBuildsRequest(projConfig.CircleCIURL, projConfig.CircleCIToken, limit, offset)
	if err != nil {
		return
	}
	err = api.NoRedirectClientDo(req, &builds)
	return
}

// FindByBuildParameters returns the first build where searched parameters have the same values
func FindByBuildParameters(builds *[]CIBuildResponse, params aws.BuildParameters) *CIBuildResponse {
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
