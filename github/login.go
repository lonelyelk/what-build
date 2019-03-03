package github

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/spf13/viper"

	"github.com/lonelyelk/what-build/api"
	"github.com/manifoldco/promptui"
	gonanoid "github.com/matoous/go-nanoid"
)

// AuthRequest is a structure for authentication request body
type AuthRequest struct {
	Note        string   `json:"note"`
	Scopes      []string `json:"scopes"`
	Fingerprint string   `json:"fingerprint"`
}

// Authorization is a body structure of response with token
type Authorization struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

const githubTokenConfigName = "github_token"

// GetAuthRequest initializes AuthRequest
func GetAuthRequest() *AuthRequest {
	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil
	}
	return &AuthRequest{
		Note:        "what-build-access",
		Scopes:      []string{"repo"},
		Fingerprint: id,
	}
}

func promptUsernameAndPassword() (username string, password string, err error) {
	prompt := promptui.Prompt{
		Label: "Enter GitHub username",
	}
	username, err = prompt.Run()
	if err != nil {
		return
	}
	prompt = promptui.Prompt{
		Label: "Enter GitHub password (not stored anywhere)",
		Mask:  '*',
	}
	password, err = prompt.Run()
	return
}

func getToken() (string, error) {
	b, _ := json.Marshal(GetAuthRequest())
	req, err := http.NewRequest("POST", "https://api.github.com/authorizations", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	username, password, err := promptUsernameAndPassword()
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)

	var a Authorization
	err = api.NoRedirectClientDo(req, &a)
	if err == nil {
		return a.Token, nil
	}
	if reflect.DeepEqual(err, api.ErrStatus(req.URL, http.StatusUnauthorized)) &&
		strings.Contains(a.Message, "OTP") {
		prompt := promptui.Prompt{
			Label: "Enter two step authentication code",
			Mask:  '*',
		}
		otp, err := prompt.Run()
		if err != nil {
			return "", err
		}
		req.Header.Add("X-GitHub-OTP", otp)
		err = api.NoRedirectClientDo(req, &a)
		return a.Token, err
	}
	return "", err
}

// Auth checks if token exists in local config and requests new one if it doesn't
func Auth() error {
	token := viper.GetString(githubTokenConfigName)
	if token != "" {
		return nil
	}
	token, err := getToken()
	if err != nil {
		return err
	}
	viper.Set(githubTokenConfigName, token)
	return viper.WriteConfig()
}
