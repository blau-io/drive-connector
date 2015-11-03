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

func AuthURL() string {
	return config.AuthCodeURL("google", oauth2.AccessTypeOffline)
}

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

func Validate(code string) (string, time.Time, error) {
	token, err := config.Exchange(oauth2.NoContext, code)
	return token.AccessToken, token.Expiry, err
}
