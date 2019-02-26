package circleci_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lonelyelk/what-build/aws"
	"github.com/lonelyelk/what-build/circleci"
)

func TestFetchBuildsRequest(t *testing.T) {
	req, err := circleci.FetchBuildsRequest("https://some.url", "some_token", 10, 20)
	if err != nil {
		t.Errorf("Expected error to be nil")
	}
	if req.URL.RawQuery != "circle-token=some_token&limit=10&offset=20" {
		t.Errorf("Expected '%s' to be 'circle-token=some_token&limit=10&offset=20'", req.URL.RawQuery)
	}
}

func TestFindByBuildParameters(t *testing.T) {
	searchParams := aws.BuildParameters{
		"some_param":  "some_value",
		"other_param": "other_value",
	}
	builds := []circleci.CIBuildResponse{
		{
			BuildNum: 0,
		},
		{
			BuildNum: 1,
			BuildParameters: aws.BuildParameters{
				"q":          2,
				"some_param": "some_value",
			},
		},
		{
			BuildNum: 2,
			BuildParameters: aws.BuildParameters{
				"other_param": "other_value",
				"some_param":  "some_value",
				"yet_another": false,
			},
		},
		{
			BuildNum: 3,
			BuildParameters: aws.BuildParameters{
				"some_param":  "some_value",
				"other_param": "other_value",
			},
		},
	}
	build := circleci.FindByBuildParameters(&builds, searchParams)
	if build == nil {
		t.Errorf("Expected a build to be found")
	}
	if build.BuildNum != 2 {
		t.Errorf("Expected build %v to be the build 2", build)
	}
}

func TestFindByBuildParameters_NotFound(t *testing.T) {
	searchParams := aws.BuildParameters{
		"some_param":  "some_value",
		"other_param": "other_value",
	}
	builds := []circleci.CIBuildResponse{
		{
			BuildNum: 1,
			BuildParameters: aws.BuildParameters{
				"q":          2,
				"some_param": "some_value",
			},
		},
	}
	build := circleci.FindByBuildParameters(&builds, searchParams)
	if build != nil {
		t.Errorf("Expected %v to be nil", build)
	}
}

func TestFindBuildsDo_Redirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://lonelyelk.ru", http.StatusFound)
	}))
	token := "tokensecret"
	prj := &aws.Project{
		Name:          "a",
		CircleCIURL:   ts.URL,
		CircleCIToken: token,
	}
	_, err := circleci.FetchBuildsDo(prj, 1, 0)
	if err == nil {
		t.Errorf("Expected redirect response to fail project fetch")
	}
	if !strings.Contains(err.Error(), ts.URL) {
		t.Errorf("Expected error '%s' to contain project url '%s'", err.Error(), ts.URL)
	}
	if strings.Contains(err.Error(), token) {
		t.Errorf("Expected error '%s' not to contain token '%s'", err.Error(), token)
	}
}

func TestFindBuildsDo_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
	}))
	token := "tokensecret"
	prj := &aws.Project{
		Name:          "a",
		CircleCIURL:   ts.URL,
		CircleCIToken: token,
	}
	_, err := circleci.FetchBuildsDo(prj, 1, 0)
	if err == nil {
		t.Errorf("Expected not found response to fail project fetch")
	}
	if !strings.Contains(err.Error(), ts.URL) {
		t.Errorf("Expected error '%s' to contain project url '%s'", err.Error(), ts.URL)
	}
	if strings.Contains(err.Error(), token) {
		t.Errorf("Expected error '%s' not to contain token '%s'", err.Error(), token)
	}
}
