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

func TestAdd(t *testing.T) {
	router := httprouter.New()
	router.POST("/add/*filepath", add)

	var addTestTable = []struct {
		file   string
		token  string
		status int
	}{
		{"", "", http.StatusUnauthorized},
		{"", "random", http.StatusBadRequest},
		{"test", "random", http.StatusOK},
	}

	for _, test := range addTestTable {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "http://x.co/add/"+test.file, nil)

		if test.token != "" {
			r.AddCookie(&http.Cookie{Name: "token", Value: test.token})
		}

		router.ServeHTTP(w, r)

		if w.Code != test.status {
			t.Errorf("Wanted Status %d, got %d", test.status, w.Code)
		}
	}
}

func TestAuthURL(t *testing.T) {
	router := httprouter.New()
	router.GET("/auth/new/:provider", authURL)

	var authURLTestTable = []struct {
		url    string
		status int
	}{
		{"/auth/new/random", http.StatusNotFound},
		{"/auth/new/google", http.StatusOK},
	}

	for _, test := range authURLTestTable {
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
		var v authURLJSON
		if err := dec.Decode(&v); err != nil {
			t.Errorf("Error while decoding json: %s", err.Error())
			continue
		}

		_, err := url.Parse(v.URL)
		if err != nil {
			t.Errorf("Error while validating URL: %s", err.Error())
			continue
		}
	}
}

func TestBrowse(t *testing.T) {
	router := httprouter.New()
	router.GET("/browse/*filepath", browse)

	var browseTestTable = []struct {
		file   string
		token  string
		status int
	}{
		{"", "", http.StatusUnauthorized},
		{"", "random", http.StatusOK},
	}

	for _, test := range browseTestTable {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://x.co/browse/"+test.file, nil)

		if test.token != "" {
			r.AddCookie(&http.Cookie{Name: "token", Value: test.token})
		}

		router.ServeHTTP(w, r)

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
		var b browseJSON
		if err := dec.Decode(&b); err != nil {
			t.Errorf("Error while decoding json: %s", err.Error())
		}
	}
}

func TestPublish(t *testing.T) {
	router := httprouter.New()
	router.GET("/publish/*filepath", publish)

	var publishTestTable = []struct {
		path   string
		token  string
		status int
	}{
		{"test", "", http.StatusUnauthorized},
		{"", "test", http.StatusBadRequest},
		{"test", "test", http.StatusOK},
	}

	for _, test := range publishTestTable {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://x.co/publish/"+test.path, nil)

		if test.token != "" {
			r.AddCookie(&http.Cookie{Name: "token", Value: test.token})
		}

		router.ServeHTTP(w, r)

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
		var v publishJSON
		if err := dec.Decode(&v); err != nil {
			t.Errorf("Error while decoding json: %s", err.Error())
		}
	}
}

func TestRead(t *testing.T) {
	router := httprouter.New()
	router.GET("/read/*filepath", read)

	var readTestTable = []struct {
		filepath string
		token    string
		status   int
	}{
		{"", "", http.StatusUnauthorized},
		{"", "test", http.StatusBadRequest},
		{"example", "test", http.StatusOK},
	}

	for _, test := range readTestTable {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://x.co/read/"+test.filepath, nil)

		if test.token != "" {
			r.AddCookie(&http.Cookie{Name: "token", Value: test.token})
		}

		router.ServeHTTP(w, r)

		if w.Code != test.status {
			t.Errorf("Wanted Status %d, got %d", test.status, w.Code)
		}
	}
}

func TestRemove(t *testing.T) {
	router := httprouter.New()
	router.DELETE("/remove/*filepath", remove)

	var removeTestTable = []struct {
		file   string
		token  string
		status int
	}{
		{"", "", http.StatusUnauthorized},
		//{"", "token", http.StatusBadRequest},
		{"test", "", http.StatusUnauthorized},
		{"test", "token", http.StatusOK},
	}

	for _, test := range removeTestTable {
		r, _ := http.NewRequest("DELETE", "http://x.co/remove/"+test.file, nil)
		w := httptest.NewRecorder()

		if test.token != "" {
			r.AddCookie(&http.Cookie{Name: "token", Value: test.token})
		}

		router.ServeHTTP(w, r)

		if w.Code != test.status {
			t.Errorf("Wanted Status %d, got %d", test.status, w.Code)
		}
	}
}

func TestValidate(t *testing.T) {
	var validateTestTable = []struct {
		formkey1   string
		formvalue1 string
		formkey2   string
		formvalue2 string
		status     int
	}{
		{"", "", "", "", http.StatusBadRequest},
		{"state", "random", "code", "test", http.StatusBadRequest},
		//{"state", "google", "code", "invalid", http.StatusBadRequest},
		{"state", "google", "code", "test", http.StatusOK},
	}

	for _, test := range validateTestTable {
		form := url.Values{}
		form.Set(test.formkey1, test.formvalue1)
		form.Set(test.formkey2, test.formvalue2)

		d := strings.NewReader(form.Encode())
		r, _ := http.NewRequest("POST", "http://x.co", d)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		validate(w, r, nil)

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
		var v validateJSON
		if err := dec.Decode(&v); err != nil {
			t.Errorf("Error while decoding json: %s", err.Error())
			continue
		}
	}
}
