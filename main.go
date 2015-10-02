package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	config      *oauth2.Config
	globalFlags struct {
		ClientSecretFile string
		Port             int
	}
)

func init() {
	flag.StringVar(&globalFlags.ClientSecretFile, "secretFile",
		"client_secret.json", "Path to the Google Drive client secret file")
	flag.IntVar(&globalFlags.Port, "port", 80, "The Port to listen on")
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

	router := httprouter.New()
	router.GET("/auth/new", NewUser)
	router.GET("/browse/:folderid", Browse)
	router.GET("/read/*filepath", Read)
	router.POST("/auth/validate", Validate)

	log.Println("Listening on port " + string(globalFlags.Port))
	http.ListenAndServe(":"+string(globalFlags.Port), router)
}
