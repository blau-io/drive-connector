package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
)

var AuthURLtable = []struct {
	url  string
	code int
}{
	{"/auth/new/random", http.StatusNotFound},
	{"/auth/new/google", http.StatusOK},
}

func TestAdd(t *testing.T) {
}

func TestAuthURL(t *testing.T) {
	var w *httptest.ResponseRecorder

	router := httprouter.New()
	router.GET("/auth/new/:provider", AuthURL)

	for _, test := range AuthURLtable {
		r, _ := http.NewRequest("GET", test.url, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, r)

		if w.Code != test.code {
			t.Errorf("Wanted Status %d, got %d", test.code, w.Code)
		}
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Wanted Content-Type application/json, got %s",
			w.Header().Get("Content-Type"))
	}

	dec := json.NewDecoder(w.Body)
	var v AuthURLjson
	if err := dec.Decode(&v); err != nil {
		t.Errorf("Error while decoding json: %s", err.Error())
	}

	parsedURI, err := url.Parse(v.URL)
	if err != nil {
		t.Errorf("Error while validating URL: %s", err.Error())
	}

	if !parsedURI.IsAbs() {
		t.Errorf("Want an absolute URL, got: %s", v.URL)
	}
}

func TestBrowse(t *testing.T) {
}

func TestPublish(t *testing.T) {
}

func TestRead(t *testing.T) {
}

func TestRemove(t *testing.T) {
}

func TestValidate(t *testing.T) {
}
