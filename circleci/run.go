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
func TriggerBuildDo(projConfig *aws.Project, buildCfg *aws.Build, branch string, extraParams aws.BuildParameters) (build CIBuildResponse, err error) {
	aws.FetchTokenIfMissing(projConfig)
	var b strings.Builder
	fmt.Fprint(&b, projConfig.CircleCIURL)
	if projConfig.CircleCIURL[len(projConfig.CircleCIURL)-1] == '/' {
		fmt.Fprint(&b, "tree/")
	} else {
		fmt.Fprint(&b, "/tree/")
	}
	fmt.Fprint(&b, branch)
	var params aws.BuildParameters
	params = make(aws.BuildParameters)
	for key, val := range buildCfg.RunBuildParameters {
		params[key] = val
	}
	for key, val := range extraParams {
		params[key] = val
	}
	req, err := TriggerBuildRequest(b.String(), projConfig.CircleCIToken, params)
	if err != nil {
		return
	}
	err = api.NoRedirectClientDo(req, &build)
	return
}

// RunBuild looks for a build in CircleCI
func RunBuild(projCfg *aws.Project, buildCfg *aws.Build, branch string, extraParams aws.BuildParameters) (*CIBuildResponse, error) {
	cib, err := TriggerBuildDo(projCfg, buildCfg, branch, extraParams)
	if err != nil {
		return nil, err
	}
	return &cib, nil
}
