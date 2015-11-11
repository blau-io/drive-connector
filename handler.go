package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err = gd.Add(cookie.Value, r.Body, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type authURLJSON struct {
	URL string `json:"url"`
}

func authURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var a = authURLJSON{}

	switch ps.ByName("provider") {
	default:
		http.Error(w, "Provider not found", http.StatusNotFound)
		return

	case "google":
		a = authURLJSON{URL: gd.AuthURL()}
	}

	j, _ := json.Marshal(a)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))
}

type browseJSON struct {
	FileList []string `json:"file_list"`
}

func browse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "authenticated", http.StatusUnauthorized)
		return
	}

	list, err := gd.Browse(cookie.Value, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	j, _ := json.Marshal(browseJSON{FileList: list})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))
}

type publishJSON struct {
	URL string `json:"url"`
}

func publish(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	link, err := gd.Publish(cookie.Value, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	j, _ := json.Marshal(publishJSON{URL: link})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))
}

func read(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	response, err := gd.Read(cookie.Value, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if response == nil {
		return
	}

	io.Copy(w, response.Body)
	response.Body.Close()
}

func remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err = gd.Delete(cookie.Value, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type validateJSON struct {
	Token  string    `json:"access_token"`
	Expiry time.Time `json:"expiry,omitempty"`
}

func validate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	state := r.FormValue("state")

	var err error
	var token string
	var expiry time.Time

	// TODO: validate state token and map to a specific provider
	switch state {
	default:
		http.Error(w, "Invalid state token", http.StatusBadRequest)
		return

	case "google":
		token, expiry, err = gd.Validate(r.FormValue("code"))
	}

	if err != nil {
		http.Error(w, "Auth Code invalid", http.StatusBadRequest)
		return
	}

	j, _ := json.Marshal(validateJSON{Token: token, Expiry: expiry})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))
}
