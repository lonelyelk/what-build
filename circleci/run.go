package circleci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lonelyelk/what-build/api"
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

	req, err = http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("circle-token", token)
	req.URL.RawQuery = q.Encode()
	return
}

// TriggerBuildDo makes a POST request to a URL from project config to trigger Circle CI build
func TriggerBuildDo(projectConfig *aws.Project, buildCfg *aws.Build, branch string) (build CIBuildResponse, err error) {
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
	var b strings.Builder
	fmt.Fprint(&b, projectConfig.CircleCIURL)
	if projectConfig.CircleCIURL[len(projectConfig.CircleCIURL)-1] == '/' {
		fmt.Fprint(&b, "tree/")
	} else {
		fmt.Fprint(&b, "/tree/")
	}
	fmt.Fprint(&b, branch)
	req, err := TriggerBuildRequest(b.String(), projectConfig.CircleCIToken, buildCfg.RunBuildParameters)
	if err != nil {
		return
	}
	err = api.NoRedirectClientDo(req, &build)
	return
}

// RunBuild looks for a build in CircleCI
func RunBuild(projCfg *aws.Project, buildCfg *aws.Build, branch string) (*CIBuildResponse, error) {
	cib, err := TriggerBuildDo(projCfg, buildCfg, branch)
	if err != nil {
		return nil, err
	}
	return &cib, nil
}
