package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/blau-io/warehouse-manager/warehouse"
	"github.com/julienschmidt/httprouter"
)

var (
	flags struct {
		GoogleSecretFile string
		Port             int
	}
	wh *warehouse.Client
)

func init() {
	flag.StringVar(&flags.GoogleSecretFile, "secretFile",
		"client_secret.json", "Path to the Google Drive client secret file")
	flag.IntVar(&flags.Port, "port", 80, "The Port to listen on")
	flag.Parse()

	var err error
	wh, err = warehouse.NewClient(flags.GoogleSecretFile)
	if err != nil {
		log.Printf("Error: %s. Proceeding...", err.Error())
	}
}

func main() {
	router := httprouter.New()

	router.POST("/add/*filepath", add)
	router.GET("/auth/new/:provider", authURL)
	router.POST("/auth/validate", validate)
	router.GET("/browse/*filepath", browse)
	router.GET("/publish/*filepath", publish)
	router.GET("/read/*filepath", read)
	router.DELETE("/remove/*filepath", remove)

	log.Printf("Listening on port %d\n", flags.Port)
	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(flags.Port), router))
}
