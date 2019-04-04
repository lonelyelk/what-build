package github

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/spf13/viper"

	"github.com/lonelyelk/what-build/api"
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

// GHBranchResponse is a structure of GitHub response
type GHBranchResponse struct {
	Name string `json:"name"`
}

const (
	developBranchName = "develop"
	branchesOption    = "...other (choose a branch directly)"
	prsOption         = "...back to PRs"
)

// ListRequest constructs and returns GitHub API based request for listing pull requests
func ListRequest(url string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("access_token", viper.GetString(githubTokenConfigName))
	req.URL.RawQuery = q.Encode()
	return
}

// ListPRsDo makes a GET request to get app open PRs
func ListPRsDo(projectConfig *aws.Project) (prs []GHPRResponse, err error) {
	req, err := ListRequest(projectConfig.GitHubURL)
	if err != nil {
		return
	}
	err = api.NoRedirectClientDo(req, &prs)
	return
}

// ListBranchesDo makes a GET request to get all remove branches
func ListBranchesDo(projectConfig *aws.Project) (prs []GHBranchResponse, err error) {
	re := regexp.MustCompile("pulls/?$")
	req, err := ListRequest(re.ReplaceAllString(projectConfig.GitHubURL, "branches"))
	if err != nil {
		return
	}
	err = api.NoRedirectClientDo(req, &prs)
	return
}

// ListAndPromptBranch lists all PRs and prompts to return a branch name
func ListAndPromptBranch(projCfg *aws.Project) (string, error) {
	prs, err := ListPRsDo(projCfg)
	if err != nil {
		return "", err
	}
	names := make([]string, len(prs)+2)
	names[0] = developBranchName
	for i, pr := range prs {
		names[i+1] = pr.Title
	}
	names[len(prs)+1] = branchesOption

	prompt := promptui.Select{
		Label: "Select a PR or branch",
		Items: names,
	}
	_, name, err := prompt.Run()

	if err != nil {
		return "", err
	}
	if name == developBranchName {
		return name, nil
	}
	if name == branchesOption {
		return ListAndPromptBranchByName(projCfg)
	}
	for _, pr := range prs {
		if pr.Title == name {
			return pr.Head.Ref, nil
		}
	}
	return "", errors.New("Not found")
}

// ListAndPromptBranchByName lists all remote branches and prompts to return a branch name
func ListAndPromptBranchByName(projCfg *aws.Project) (string, error) {
	prs, err := ListBranchesDo(projCfg)
	if err != nil {
		return "", err
	}
	names := make([]string, len(prs)+1)
	names[0] = prsOption
	for i, pr := range prs {
		names[i+1] = pr.Name
	}

	prompt := promptui.Select{
		Label: "Select a branch",
		Items: names,
	}
	_, name, err := prompt.Run()

	if err != nil {
		return "", err
	}
	if name == prsOption {
		return ListAndPromptBranch(projCfg)
	}
	return name, nil
}
