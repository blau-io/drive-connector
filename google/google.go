package google

import (
	"io/ioutil"

	"golang.org/x/oauth2"
	goauth "golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	config *oauth2.Config
)

func Config(filepath string) error {
	secret, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	config, err = goauth.ConfigFromJSON(secret, drive.DriveScope)
	if err != nil {
		return err
	}

	return nil
}
