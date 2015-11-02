package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AuthUrlJSON struct {
	URL string
}

func Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func AuthUrl(w http.ResponseWriter, r *http.Request,
	ps httprouter.Params) {
	var a = AuthUrlJSON{}

	switch ps.ByName("provider") {
	default:
		http.Error(w, "Provider not found", http.StatusNotFound)
		return

	case "google":
		a = AuthUrlJSON{URL: "http://google.com"}
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
