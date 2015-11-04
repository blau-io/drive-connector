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
	flag.StringVar(&globalFlags.GoogleSecretFile, "secretFile", "",
		"Path to the Google Drive client secret file")
	flag.IntVar(&globalFlags.Port, "port", 80, "The Port to listen on")
	flag.Parse()

	if globalFlags.GoogleSecretFile != "" {
		if err := googledrive.Config(globalFlags.GoogleSecretFile); err != nil {
			log.Fatalf("Could not configure Google integration: %s", err.Error())
		}
	}
}

func main() {
	router := httprouter.New()

	// Add
	router.POST("/auth/new/*filepath", Add)

	// Auth
	router.GET("/auth/new/:provider", AuthURL)
	router.POST("/auth/validate", Validate)

	log.Println("Listening on port " + strconv.Itoa(globalFlags.Port))
	log.Println(http.ListenAndServe(":"+strconv.Itoa(globalFlags.Port), router))
}
