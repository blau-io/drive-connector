package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"golang.org/x/oauth2"
)

func auth(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		url := config.AuthCodeURL("todo", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Unable to retrieve token: %v", err)
		http.Error(w, "Unable to retrieve token",
			http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, token.AccessToken)
}

func index(w http.ResponseWriter, r *http.Request) {
	srv, err := getClient(r)
	if err != nil {
		log.Printf("Unable to retrieve drive Client: %v", err)
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}

	file, err := getFileByPath(srv, path.Join("/blau.io/configuration/",
		r.URL.Path))
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		http.Error(w, "Unable to retrieve files", http.StatusUnauthorized)
		return
	}

	response, err := srv.Files.Get(file.Id).Download()
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	io.Copy(w, response.Body)
	response.Body.Close()
}

func list(w http.ResponseWriter, r *http.Request) {
	srv, err := getClient(r)
	if err != nil {
		log.Printf("Unable to retrieve drive Client: %v", err)
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}

	url := strings.TrimPrefix(r.URL.Path, "/list")
	list, err := getDirectoryList(srv, path.Join("/blau.io/content/", url))
	if err != nil {
		log.Printf("Unable to get list of directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonResponse := make([]string, len(list))
	for i, v := range list {
		jsonResponse[i] = v.Title
	}

	j, err := json.Marshal(jsonResponse)
	if err != nil {
		log.Printf("Error while encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprint(w, string(j))
}
