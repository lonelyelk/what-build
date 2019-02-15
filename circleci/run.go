package circleci

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lonelyelk/what-build/aws"
)

// PostParams is a struct for Circle CI POST request params
type PostParams struct {
	BuildParameters aws.BuildParameters `json:"build_parameters"`
}

// TriggerBuildRequest constructs and returns CircleCI API based request for triggering builds
func TriggerBuildRequest(url string, token string, params aws.BuildParameters) (req *http.Request, err error) {
	params["IAM_USER"] = aws.GetIAMUserName()
	postParams := PostParams{BuildParameters: params}
	b, err := json.Marshal(postParams)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest("POST", url+"/tree/develop", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("circle-token", token)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "lonelyelk-what-build-bot")
	return
}

// TriggerBuildDo makes a POST request to a URL from project config to trigger Circle CI build
func TriggerBuildDo(projectConfig *aws.Project, buildCfg *aws.Build) (build CIBuildResponse, err error) {
	if projectConfig.CircleCIToken == "" {
		var token string
		token, err = aws.GetSSMParameter(projectConfig.CircleCITokenSSMName)
		if err != nil {
			return
		}
		// Account for post request prepared token like 'token:' as if it was user with no password
		if token[len(token)-1] == ':' {
			token = token[:len(token)-1]
		}
		projectConfig.CircleCIToken = token
	}
	req, err := TriggerBuildRequest(projectConfig.CircleCIURL, projectConfig.CircleCIToken, buildCfg.RunBuildParameters)
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
	if res.StatusCode != 201 {
		return build, errStatus(req.URL)
	}
	err = json.NewDecoder(res.Body).Decode(&build)
	return
}

// RunBuild looks for a build in CircleCI
func RunBuild(projCfg *aws.Project, buildCfg *aws.Build) (*CIBuildResponse, error) {
	cib, err := TriggerBuildDo(projCfg, buildCfg)
	if err != nil {
		return nil, err
	}
	return &cib, nil
}
