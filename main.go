package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/blau-io/warehouse-manager/googledrive"
	"github.com/julienschmidt/httprouter"
)

var (
	flags struct {
		GoogleSecretFile string
		Port             int
	}
	gd *googledrive.GoogleDrive
)

func init() {
	flag.StringVar(&flags.GoogleSecretFile, "secretFile",
		"client_secret.json", "Path to the Google Drive client secret file")
	flag.IntVar(&flags.Port, "port", 80, "The Port to listen on")
	flag.Parse()

	var err error
	gd, err = googledrive.NewGoogleDrive(flags.GoogleSecretFile)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		log.Println("Google Drive not configured, proceeding...")
		gd = &googledrive.GoogleDrive{}
	}
}

func main() {
	router := httprouter.New()

	router.POST("/add/*filepath", Add)
	router.GET("/auth/new/:provider", AuthURL)
	router.POST("/auth/validate", Validate)

	log.Printf("Listening on port %d\n", flags.Port)
	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(flags.Port), router))
}
