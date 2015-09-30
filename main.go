package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"golang.org/x/net/context"
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

func auth(w http.ResponseWriter, r *http.Request) {
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

func index(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Printf("Unable to read cookie: %v", err)
		http.Error(w, "Unable to read cookie", http.StatusUnauthorized)
		return
	}

	token := &oauth2.Token{
		AccessToken: cookie.Value,
	}

	srv, err := drive.New(config.Client(context.Background(), token))
	if err != nil {
		log.Printf("Unable to retrieve drive Client: %v", err)
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}

	file, err := getFileByPath(srv, path.Join("/blau.io/configuration/"+
		r.URL.Path))
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		http.Error(w, "Unable to retrieve files", http.StatusUnauthorized)
		return
	}

	response, err := srv.Files.Get(file.Id).Download()
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	io.Copy(w, response.Body)
	response.Body.Close()
}

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
	http.ListenAndServe(":"+globalFlags.Port, nil)
}
