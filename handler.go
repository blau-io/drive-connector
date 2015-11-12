package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/blau-io/warehouse-manager/warehouse"
	"github.com/julienschmidt/httprouter"
)

func add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err = wh.Add(cookie.Value, ps.ByName("filepath"), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func authURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url, err := wh.AuthURL(ps.ByName("provider"))
	if err != nil {
		http.Error(w, "Provider not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(warehouse.AuthURLJSON{URL: url})
	w.Write([]byte(j))
}

func browse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	list, err := wh.Browse(cookie.Value, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(warehouse.BrowseJSON{FileList: list})
	w.Write([]byte(j))
}

func publish(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	link, err := wh.Publish(cookie.Value, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(warehouse.PublishJSON{URL: link})
	w.Write([]byte(j))
}

func read(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	response, err := wh.Read(cookie.Value, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if response == nil {
		// Empty file
		return
	}

	io.Copy(w, response)
	response.Close()
}

func remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err = wh.Remove(cookie.Value, ps.ByName("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func validate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token, expiry, err := wh.Validate(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Auth Code invalid", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(warehouse.ValidateJSON{Token: token, Expiry: expiry})
	w.Write([]byte(j))
}
