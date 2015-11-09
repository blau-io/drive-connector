package googledrive

import (
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v2"
)

func getClient(config *oauth2.Config, code string) (*drive.Service, error) {
	if config == nil {
		return nil, errors.New("Google Drive is not configured")
	}

	token := &oauth2.Token{AccessToken: code}
	return drive.New(config.Client(context.Background(), token))
}
