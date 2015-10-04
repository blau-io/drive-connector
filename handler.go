package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v2"
)

func Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	srv, err := getClient(r)
	if err != nil {
		log.Printf("Unable to retrieve drive Client: %v", err)
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}

	parent, err := getParent(srv, ps.ByName("filepath"))
	if err != nil {
		log.Printf("Unable to get parent id: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	paths := strings.Split(strings.TrimPrefix(ps.ByName("filepath"), "/"), "/")

	file := &drive.File{
		Title:   paths[len(paths)-1],
		Parents: []*drive.ParentReference{parent},
	}

	if r.Header.Get("folder") == "true" {
		file.MimeType = "application/vnd.google-apps.folder"
	}

	_, err = srv.Files.Insert(file).Media(r.Body).Do()
	if err != nil {
		log.Printf("Error while uploading file: %v", err)
		http.Error(w, "Error while uploading file",
			http.StatusInternalServerError)
	}

	r.Body.Close()
}

func Browse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	srv, err := getClient(r)
	if err != nil {
		log.Printf("Unable to retrieve drive Client: %v", err)
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}

	list, err := getDirectoryList(srv, ps.ByName("filepath"))
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

func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err = srv.Files.Delete(file.Id).Do()
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func NewUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	url := config.AuthCodeURL("todo", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func Publish(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	srv, err := getClient(r)
	if err != nil {
		log.Printf("Unable to retrieve drive Client: %v", err)
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}

	file, err := getFileByPath(srv, ps.ByName("filepath"))
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		http.Error(w, "Unable to retrieve files", http.StatusNotFound)
		return
	}

	permission := &drive.Permission{
		Value: "",
		Type:  "anyone",
		Role:  "reader",
	}

	_, err = drive.NewPermissionsService(srv).Insert(file.Id, permission).Do()
	if err != nil {
		log.Printf("Unable to insert new permission: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	updated, err := srv.Files.Get(file.Id).Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		http.Error(w, "Unable to retrieve files", http.StatusNotFound)
		return
	}

	fmt.Fprint(w, updated.WebViewLink)
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
