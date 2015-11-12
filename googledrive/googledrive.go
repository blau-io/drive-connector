package googledrive

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

// GoogleDrive holds the configuration for the Google Drive SDK
type GoogleDrive struct {
	config *oauth2.Config
}

// Add inserts a new file on Google Drive
func (d *GoogleDrive) Add(code, filepath string, content io.ReadCloser) error {
	if filepath == "" || filepath == "/" {
		return errors.New("No filepath specified")
	}

	client, _ := getClient(d.config, code)
	if client == nil {
		return nil
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

// Browse return the content of a directory as a list
func (d *GoogleDrive) Browse(code, filepath string) ([]string, error) {
	client, _ := getClient(d.config, code)
	if client == nil {
		return nil, nil
	}

	folder, err := getFileByPath(client, filepath)
	if err != nil {
		return nil, err
	}

	query := "'" + folder.Id + "' in parents and trashed = false"
	list, err := client.Files.List().Q(query).Do()
	if err != nil {
		return nil, err
	}

	out := make([]string, len(list.Items))
	for i, v := range list.Items {
		out[i] = v.Title
	}

	return out, nil
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

// Publish sets a file to public and returns its public url
func (d *GoogleDrive) Publish(code, filepath string) (string, error) {
	if filepath == "" || filepath == "/" {
		return "", errors.New("Can't publish root folder")
	}

	client, _ := getClient(d.config, code)
	if client == nil {
		return "", nil
	}

	file, err := getFileByPath(client, filepath)
	if err != nil {
		return "", err
	}

	perm := &drive.Permission{
		Value: "",
		Type:  "anyone",
		Role:  "reader",
	}

	_, err = drive.NewPermissionsService(client).Insert(file.Id, perm).Do()
	if err != nil {
		return "", err
	}

	// Get the updated file, which will now contain a public URL
	file, err = client.Files.Get(file.Id).Do()
	if err != nil {
		return "", err
	}

	return file.WebViewLink, nil
}

// Read returns the content of a file
func (d *GoogleDrive) Read(code, filepath string) (*http.Response, error) {
	if filepath == "" || filepath == "/" {
		return nil, errors.New("For the content of a folder, please use browse")
	}

	client, _ := getClient(d.config, code)
	if client == nil {
		return nil, nil
	}

	file, err := getFileByPath(client, filepath)
	if err != nil {
		return nil, err
	}

	return client.Files.Get(file.Id).Download()
}

// Remove deletes a file from Google Drive
func (d *GoogleDrive) Remove(code, filepath string) error {
	if filepath == "" || filepath == "/" {
		return errors.New("Deleting the root folder is not permitted")
	}

	client, _ := getClient(d.config, code)
	if client == nil {
		return nil
	}

	file, err := getFileByPath(client, filepath)
	if err != nil {
		return err
	}

	err = client.Files.Delete(file.Id).Do()
	if err != nil {
		return err
	}

	return nil
}

// Validate validates an access code against the oauth2.config object. It
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
