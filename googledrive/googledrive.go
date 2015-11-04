package googledrive

import (
	"errors"
	"io"
	"io/ioutil"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	config     *oauth2.Config
	configured = false
)

// Add inserts a new file on Google Drive
func Add(code string, body io.ReadCloser, filepath string) error {
	if !configured {
		return errors.New("Google Drive is not configured")
	}

	token := &oauth2.Token{
		AccessToken: code,
	}

	_, err := drive.New(config.Client(context.Background(), token))
	if err != nil {
		return err
	}

	return nil
}

// AuthURL returns a URL to the Google OAuth2 login page
func AuthURL() string {
	if !configured {
		return ""
	}

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

	configured = true
	return nil
}

// Validate validates an access code against the oauth2.config object. It
// then returns the real token togehter with an expiry date.
func Validate(code string) (string, time.Time, error) {
	if !configured {
		return "", time.Now(), errors.New("Google Drive is not configured")
	}

	token, err := config.Exchange(oauth2.NoContext, code)
	return token.AccessToken, token.Expiry, err
}
