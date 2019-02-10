package circleci_test

import (
	"testing"

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
	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected request to have 'Accept application/json' header")
	}
}

func TestFindByBuildParameters(t *testing.T) {
	searchParams := map[string]interface{}{
		"some_param":  "some_value",
		"other_param": "other_value",
	}
	builds := []circleci.CIBuildResponse{
		{
			BuildNum: 0,
		},
		{
			BuildNum: 1,
			BuildParameters: map[string]interface{}{
				"q":          2,
				"some_param": "some_value",
			},
		},
		{
			BuildNum: 2,
			BuildParameters: map[string]interface{}{
				"other_param": "other_value",
				"some_param":  "some_value",
				"yet_another": false,
			},
		},
		{
			BuildNum: 3,
			BuildParameters: map[string]interface{}{
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
	searchParams := map[string]interface{}{
		"some_param":  "some_value",
		"other_param": "other_value",
	}
	builds := []circleci.CIBuildResponse{
		{
			BuildNum: 1,
			BuildParameters: map[string]interface{}{
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