package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

func Browse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	srv, err := getClient(r)
	if err != nil {
		log.Printf("Unable to retrieve drive Client: %v", err)
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}

	list, err := getDirectoryList(srv, ps.ByName("folderid"))
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
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(j))
}

func NewUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	url := config.AuthCodeURL("todo", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func Read(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	srv, err := getClient(r)
	if err != nil {
		log.Printf("Unable to retrieve drive Client: %v", err)
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}

	file, err := getFileByPath(srv, ps.ByName("filepath"))
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

func Validate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	code := r.FormValue("code")
	if code == "" {
		log.Println("No authorization code present")
		http.Error(w, "No authorization code present", http.StatusUnauthorized)
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
