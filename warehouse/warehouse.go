package warehouse

import (
	"errors"
	"io"
	"time"

	gd "github.com/blau-io/warehouse-manager/googledrive"
)

// Client holds the configuration for all providers
type Client struct {
	googledrive *gd.GoogleDrive
}

// Add adds a new file to a storage provider
func (c *Client) Add(code, filepath string, content io.ReadCloser) error {
	// TODO: Decide, based on the code, which provider to choose
	return c.googledrive.Add(code, filepath, content)
}

// AuthURL returns the url to an oauth2 providers login page
func (c *Client) AuthURL(provider string) (string, error) {
	switch provider {
	default:
		return "", errors.New("No provider specified")

	case "google":
		return c.googledrive.AuthURL(), nil
	}
}

// Browse returns a list of filenames in a folder
func (c *Client) Browse(code, filepath string) ([]string, error) {
	// TODO: Decide, based on the code, which provider to choose
	return c.googledrive.Browse(code, filepath)
}

// NewClient initializes a new warehouse client and returns it
func NewClient(googleJSON string) (*Client, error) {
	gdrive, err := gd.NewGoogleDrive(googleJSON)
	if err != nil {
		gdrive = &gd.GoogleDrive{}
	}

	return &Client{googledrive: gdrive}, err
}

// Publish makes a private file public
func (c *Client) Publish(code, filepath string) (string, error) {
	// TODO: Decide, based on the code, which provider to choose
	return c.googledrive.Publish(code, filepath)
}

// Read returns the content of a given file
func (c *Client) Read(code, filepath string) (io.ReadCloser, error) {
	response, err := c.googledrive.Read(code, filepath)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	return response.Body, nil
}

// Remove deletes a file from a storage provider
func (c *Client) Remove(code, filepath string) error {
	// TODO: Decide, based on the code, which provider to choose
	return c.googledrive.Remove(code, filepath)
}

// Validate validates an oauth2 code and returns a token
func (c *Client) Validate(state, code string) (string, time.Time, error) {
	// TODO: Validate state token and map to a specific provider
	switch state {
	default:
		return "", time.Now(), errors.New("Invalid state token")

	case "google":
		return c.googledrive.Validate(code)
	}
}
