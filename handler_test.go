package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestAdd(t *testing.T) {
}

func TestAuthUrl(t *testing.T) {
	router := httprouter.New()
	router.GET("/auth/new/:provider", AuthUrl)

	r, _ := http.NewRequest("GET", "/auth/new/random", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Wanted Status 404, got %d", w.Code)
	}

	r, _ = http.NewRequest("GET", "/auth/new/google", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Wanted Status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Wanted Content-Type application/json, got %s",
			w.Header().Get("Content-Type"))
	}

	dec := json.NewDecoder(w.Body)
	var v AuthUrlJSON
	if err := dec.Decode(&v); err != nil {
		t.Errorf("Error while decoding json: %s", err.Error())
	}

	parsedUri, err := url.Parse(v.URL)
	if err != nil {
		t.Errorf("Error while validating URL: %s", err.Error())
	}

	if !parsedUri.IsAbs() {
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
