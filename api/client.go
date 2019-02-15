package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func errStatus(url *url.URL) error {
	return fmt.Errorf("api: url '%s://%s%s' doesn't succeed", url.Scheme, url.Host, url.Path)
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
			return errStatus(r.URL)
		},
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return errStatus(req.URL)
	}
	err = json.NewDecoder(res.Body).Decode(decoder)
	return
}
