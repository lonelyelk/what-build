package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/lonelyelk/what-build/aws"
	"github.com/manifoldco/promptui"
)

// Head refers to git head of PR's branch
type Head struct {
	Ref string `json:"ref"`
}

// GHPRResponse is a structure of GitHub response
type GHPRResponse struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Head  `json:"head"`
}

func errStatus(url *url.URL) error {
	return fmt.Errorf("circleci: project url '%s://%s%s' doesn't exist", url.Scheme, url.Host, url.Path)
}
func errBuildNotFound(buildName string, projectName string) error {
	return fmt.Errorf("circleci: build '%s' not found for project '%s'", buildName, projectName)
}

// ListPRsRequest constructs and returns GitHub API based request for listing pull requests
func ListPRsRequest(url string, token string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("access_token", token)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lonelyelk-what-build-bot")
	return
}

// ListPRsDo makes a POST request to a URL from project config to trigger Circle CI build
func ListPRsDo(projectConfig *aws.Project) (prs []GHPRResponse, err error) {
	req, err := ListPRsRequest(projectConfig.GitHubURL, projectConfig.GitHubToken)
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
		return prs, errStatus(req.URL)
	}
	err = json.NewDecoder(res.Body).Decode(&prs)
	return
}

// ListAndPromptBranch lists all PRs and prompts to return a branch name
func ListAndPromptBranch(projCfg *aws.Project) (string, error) {
	prs, err := ListPRsDo(projCfg)
	if err != nil {
		return "", err
	}
	names := make([]string, len(prs)+1)
	names[0] = "develop"
	for i, pr := range prs {
		names[i+1] = pr.Title
	}

	prompt := promptui.Select{
		Label: "Select a PR or branch",
		Items: names,
	}
	_, name, err := prompt.Run()

	if err != nil {
		return "", err
	}
	if name == "develop" {
		return name, nil
	}
	for _, pr := range prs {
		if pr.Title == name {
			return pr.Head.Ref, nil
		}
	}
	return "", errors.New("Not found")
}
