package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/blau-io/warehouse-manager/googledrive"
	"github.com/julienschmidt/httprouter"
)

// Add adds a new file to a the cloud storage provider listed in the cookie
func Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Invalid oauth token", http.StatusUnauthorized)
		return
	}

	err = googledrive.Add(cookie.Value, r.Body, ps.ByName("filename"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// AuthURLJSON is the struct which will be encoded into JSON once it's been
// initialized by AuthURL().
type AuthURLJSON struct {
	URL string `json:"url"`
}

// AuthURL gets an oauth2 URL from one of the supported libraries (depending
// on httprouter.Params) and returns the link encoded in JSON.
// If httprouter.Params specify an unsupported library, http.StatusNotFound
// is returned.
func AuthURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var a = AuthURLJSON{}

	switch ps.ByName("provider") {
	default:
		http.Error(w, "Provider not found", http.StatusNotFound)
		return

	case "google":
		a = AuthURLJSON{URL: googledrive.AuthURL()}
	}

	j, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))
}

// ValidateJSON is the struct which will be encoded into JSON once it's been
// initialized by Validate().
type ValidateJSON struct {
	Token  string    `json:"access_token"`
	Expiry time.Time `json:"expiry,omitempty"`
}

// Validate reads the Form Values of a request and validates the oauth2.
// After the code is validated, it returns the user token.
func Validate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	state := r.FormValue("state")

	// TODO: validate state token and map to a specific provider
	if state != "google" {
		http.Error(w, "Invalid state token", http.StatusBadRequest)
		return
	}

	token, expiry, err := googledrive.Validate(r.FormValue("code"))
	if err != nil {
		http.Error(w, "Auth Code invalid", http.StatusBadRequest)
		return
	}

	v := ValidateJSON{Token: token, Expiry: expiry}

	j, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))
}
