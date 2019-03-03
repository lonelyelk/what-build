package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ErrStatus formats error with url and status
func ErrStatus(url *url.URL, status int) error {
	return fmt.Errorf("api: url '%s://%s%s' returned %d", url.Scheme, url.Host, url.Path, status)
}

// NoRedirectClientDo performes http request and decodes json
func NoRedirectClientDo(req *http.Request, decoder interface{}) (err error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "lonelyelk-what-build-bot")
	client := http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			var r = req
			if len(via) > 0 {
				r = via[len(via)-1]
			}
			return ErrStatus(r.URL, http.StatusFound)
		},
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		err = ErrStatus(req.URL, res.StatusCode)
	}
	jsonErr := json.NewDecoder(res.Body).Decode(decoder)
	if err == nil {
		return jsonErr
	}
	return
}
