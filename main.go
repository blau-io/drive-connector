package main

import (
	"flag"
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
		Port             string
	}
)

func init() {
	flag.StringVar(&globalFlags.ClientSecretFile, "secretFile",
		"client_secret.json", "Path to the Google Drive client secret file")
	flag.StringVar(&globalFlags.Port, "port", "80", "The Port to listen on")
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

	log.Println("Listening on port " + globalFlags.Port)

	http.HandleFunc("/", index)
	http.HandleFunc("/auth", auth)
	http.HandleFunc("/list/", list)
	http.ListenAndServe(":"+globalFlags.Port, nil)
}
