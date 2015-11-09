package googledrive

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

// GoogleDrive holds the configuration for the Google Drive SDK
type GoogleDrive struct {
	config *oauth2.Config
}

// Add inserts a new file on Google Drive
func (d *GoogleDrive) Add(code string, content io.ReadCloser,
	filepath string) error {
	if filepath == "" || filepath == "/" {
		return errors.New("No filepath specified")
	}

	if d.config == nil {
		return nil
	}

	token := &oauth2.Token{AccessToken: code}
	client, err := drive.New(d.config.Client(context.Background(), token))
	if err != nil {
		return err
	}

	parent, err := getParent(client, filepath)
	if err != nil {
		return err
	}

	paths := strings.Split(sanitize(filepath), "/")
	file := &drive.File{
		Title:   paths[len(paths)-1],
		Parents: []*drive.ParentReference{parent},
	}

	if _, err = client.Files.Insert(file).Media(content).Do(); err != nil {
		return err
	}

	content.Close()
	return nil
}

// AuthURL returns a URL to the Google OAuth2 login page
func (d *GoogleDrive) AuthURL() string {
	if d.config == nil {
		return ""
	}

	return d.config.AuthCodeURL("google", oauth2.AccessTypeOffline)
}

// NewGoogleDrive reads the information from the supplied secret file and
// parses it into the config of a new googledrive object. It returns the newly
// created object.
func NewGoogleDrive(filepath string) (*GoogleDrive, error) {
	secret, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(secret, drive.DriveScope)
	if err != nil {
		return nil, err
	}

	return &GoogleDrive{config: config}, nil
}

// validate validates an access code against the oauth2.config object. It
// then returns the real token togehter with an expiry date.
func (d *GoogleDrive) Validate(code string) (string, time.Time, error) {
	if d.config == nil {
		return "", time.Now(), nil
	}

	token, err := d.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", time.Now(), err
	}

	return token.AccessToken, token.Expiry, nil
}
