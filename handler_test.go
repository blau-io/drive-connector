package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

var authURLtable = []struct {
	url    string
	status int
}{
	{"/auth/new/random", http.StatusNotFound},
	//	{"/auth/new/google", http.StatusOK},
}

func TestAuthURL(t *testing.T) {
	router := httprouter.New()
	router.GET("/auth/new/:provider", AuthURL)

	for _, test := range authURLtable {
		r, _ := http.NewRequest("GET", test.url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)

		if w.Code != test.status {
			t.Errorf("Wanted Status %d, got %d", test.status, w.Code)
			continue
		}

		if w.Code != http.StatusOK {
			continue
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Wanted Content-Type application/json, got %s",
				w.Header().Get("Content-Type"))
		}

		dec := json.NewDecoder(w.Body)
		var v AuthURLJSON
		if err := dec.Decode(&v); err != nil {
			t.Errorf("Error while decoding json: %s", err.Error())
			continue
		}

		parsedURI, err := url.Parse(v.URL)
		if err != nil {
			t.Errorf("Error while validating URL: %s", err.Error())
			continue
		}

		if !parsedURI.IsAbs() {
			t.Errorf("Want an absolute URL, got: %s", v.URL)
		}
	}
}

var validateTable = []struct {
	formkey1   string
	formvalue1 string
	formkey2   string
	formvalue2 string
	status     int
}{
	{"", "", "", "", http.StatusBadRequest},
	{"state", "random", "code", "test", http.StatusBadRequest},
	//	{"state", "google", "code", "test", http.StatusOK},
}

func TestValidate(t *testing.T) {
	for _, test := range validateTable {
		form := url.Values{}
		form.Set(test.formkey1, test.formvalue1)
		form.Set(test.formkey2, test.formvalue2)

		d := strings.NewReader(form.Encode())
		r, _ := http.NewRequest("POST", "http://foo.bar", d)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		Validate(w, r, nil)

		if w.Code != test.status {
			t.Errorf("Wanted Status %d, got %d", test.status, w.Code)
		}

		if w.Code != http.StatusOK {
			continue
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Wanted Content-Type application/json, got %s",
				w.Header().Get("Content-Type"))
		}

		dec := json.NewDecoder(w.Body)
		var v ValidateJSON
		if err := dec.Decode(&v); err != nil {
			t.Errorf("Error while decoding json: %s", err.Error())
			continue
		}

		if v.Token == "" {
			t.Error("Got an empty token")
		}
	}
}
