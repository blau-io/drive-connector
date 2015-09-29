package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	config      *oauth2.Config
	globalFlags struct {
		ClientSecretFile string
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
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

func init() {
	flag.StringVar(&globalFlags.ClientSecretFile, "secretFile",
		"client_secret.json", "Path to the Google Drive client secret file")
	flag.Parse()
}

func main() {
	secret, err := ioutil.ReadFile(globalFlags.ClientSecretFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err = google.ConfigFromJSON(secret, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	http.HandleFunc("/auth", handler)
	http.ListenAndServe(":8080", nil)
}
