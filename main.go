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
	globalFlags struct {
		GoogleSecretFile string
		Port             int
	}
)

func init() {
	flag.StringVar(&globalFlags.GoogleSecretFile, "secretFile",
		"client_secret.json", "Path to the Google Drive client secret file")
	flag.IntVar(&globalFlags.Port, "port", 80, "The Port to listen on")
	flag.Parse()
}

func main() {
	if err := googledrive.Config(globalFlags.GoogleSecretFile); err != nil {
		log.Fatalf("Could not configure Google integration: %s", err.Error())
	}

	router := httprouter.New()

	// Browse
	router.GET("/browse/*filepath", Browse)

	// Read
	router.GET("/read/*filepath", Read)

	// Remove
	router.DELETE("/remove/*filepath", Remove)

	// Auth
	router.GET("/auth/new/:provider", AuthURL)
	router.POST("/auth/validate", Validate)

	// Add (Create)
	router.POST("/add/*filepath", Add)

	// Publish
	router.GET("/publish/*filepath", Publish)

	log.Println("Listening on port " + strconv.Itoa(globalFlags.Port))
	log.Println(http.ListenAndServe(":"+strconv.Itoa(globalFlags.Port), router))
}
