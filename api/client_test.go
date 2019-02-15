package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lonelyelk/what-build/api"
)

type testJSON struct {
	Field string `json:"field"`
}

func TestNoRedirectClientDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected request to have 'Accept application/json' header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected request to have 'Content-Type application/json' header")
		}
		fmt.Fprint(w, `{"field": "value"}`)
	}))
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Expected to create a new request")
	}
	var d testJSON
	err = api.NoRedirectClientDo(req, &d)
	if err != nil {
		t.Errorf("Expected response not to fail with %s", err)
	}
	if d.Field != "value" {
		t.Errorf("Expected field '%s' to equal 'value'", d.Field)
	}
}

func TestNoRedirectClientDo_Redirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://lonelyelk.ru", http.StatusFound)
	}))
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Expected to create a new request")
	}
	var d interface{}
	err = api.NoRedirectClientDo(req, d)
	if err == nil {
		t.Errorf("Expected redirect response to fail fetch")
	}
	if !strings.Contains(err.Error(), ts.URL) {
		t.Errorf("Expected error '%s' to contain url '%s'", err.Error(), ts.URL)
	}
}

func TestNoRedirectClientDo_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
	}))
	req, err := http.NewRequest("POST", ts.URL, nil)
	if err != nil {
		t.Errorf("Expected to create a new request")
	}
	var d interface{}
	err = api.NoRedirectClientDo(req, d)
	if err == nil {
		t.Errorf("Expected not found response to fail fetch")
	}
	if !strings.Contains(err.Error(), ts.URL) {
		t.Errorf("Expected error '%s' to contain url '%s'", err.Error(), ts.URL)
	}
}
