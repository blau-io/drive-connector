package google

import (
	"golang.org/x/oauth2"
	goauth "golang.org/x/oauth2/google"
)

var (
	config *oauth2.Config
)

func Config(string filepath) error {
	secret, err := ioutil.ReadFile(globalFlags.ClientSecretFile)
	if err != nil {
		return err
	}

	config, err = goauth.ConfigFromJSON(secret, drive.DriveScope)
	if err != nil {
		return err
	}

	return nil
}
