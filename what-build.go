package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

type buildParameters struct {
	DeployToQA string `json:"DEPLOY_TO_QA"`
	QAEnv      string `json:"QA_ENV"`
}
type buildResponse struct {
	BuildNum        int             `json:"build_num"`
	Branch          string          `json:"branch"`
	VcsRevision     string          `json:"vcs_revision"`
	Subject         string          `json:"subject"`
	Why             string          `json:"why"`
	DontBuild       string          `json:"dont_build"`
	StopTime        string          `json:"stop_time"`
	BuildTimeMillis int             `json:"build_time_millis"`
	Status          string          `json:"status"`
	BuildParameters buildParameters `json:"build_parameters"`
}

const perPage int = 10

func main() {
	client := &http.Client{Timeout: 10 * time.Second}
	for offset := 0; offset < 100; offset = offset + 10 {
		req, err := http.NewRequest("GET", os.Getenv("PROJECT_URL"), nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		q := req.URL.Query()
		q.Add("circle-token", os.Getenv("CIRCLECI_TOKEN"))
		q.Add("limit", string(perPage))
		q.Add("offset", string(offset))
		q.Add("filter", "completed")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("User-Agent", "lonelyelk-what-build-bot")
		req.URL.RawQuery = q.Encode()
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		defer res.Body.Close()
		builds := make([]buildResponse, 0)
		err = json.NewDecoder(res.Body).Decode(&builds)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		for _, build := range builds {
			if build.Why == "api" && build.BuildParameters.DeployToQA == "true" && build.Status == "success" {
				time, err := time.Parse(time.RFC3339, build.StopTime)
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("Deployed to %s at %v\nCommit: %s\nRevision: %s\n", build.BuildParameters.QAEnv, time, build.Subject, build.VcsRevision)
				os.Exit(0)
			}
		}
	}
	fmt.Println("Deploy to QA not found")
	os.Exit(-1)
}
