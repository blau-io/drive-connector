package googledrive

import (
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	config *oauth2.Config
)

// AuthURL returns a URL to the Google OAuth2 login page
func AuthURL() string {
	return config.AuthCodeURL("google", oauth2.AccessTypeOffline)
}

// Config reads the information from the client_secret.json file and
// parses it into the global config object, so the other functions can
// access it.
func Config(filepath string) error {
	secret, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	config, err = google.ConfigFromJSON(secret, drive.DriveScope)
	if err != nil {
		return err
	}

	return nil
}

// Validate validates an access code against the oauth2.config object. It
// then returns the real token togehter with an expiry date.
func Validate(code string) (string, time.Time, error) {
	token, err := config.Exchange(oauth2.NoContext, code)
	return token.AccessToken, token.Expiry, err
}
