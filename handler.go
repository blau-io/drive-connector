package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// AuthURLjson is the struct which will be encoded into JSON once it's been
// initialized by AuthURL().
type AuthURLjson struct {
	URL string
}

func Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

// AuthURL gets an oauth2 URL from one of the supported libraries (depending
// on httprouter.Params) and returns the link encoded in JSON.
// If httprouter.Params specify an unsupported library, http.StatusNotFound
// is returned.
func AuthURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var a = AuthURLjson{}

	switch ps.ByName("provider") {
	default:
		http.Error(w, "Provider not found", http.StatusNotFound)
		return

	case "google":
		a = AuthURLjson{URL: "http://google.com"}
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Write([]byte(j))
}

func Browse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func Publish(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func Read(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func Remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func Validate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}
