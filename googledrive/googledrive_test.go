package googledrive

import (
	"bytes"
	"io"
	"testing"

	"golang.org/x/oauth2"
)

var g = GoogleDrive{
	config: &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		RedirectURL:  "REDIRECT_URL",
		Scopes:       []string{"scope1", "scope2"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://example.com/auth",
			TokenURL: "http://example.com/token",
		},
	},
}

type mockRC struct {
	io.Reader
}

func (mockRC) Close() error {
	return nil
}

func TestAdd(t *testing.T) {
	var addTestTable = []struct {
		code     string
		filepath string
		fail     bool
	}{
		{"", "", true},
		{"", "test", true},
		{"token", "", true},
		{"token", "test", false},
	}

	for _, test := range addTestTable {
		err := g.Add(test.code, mockRC{bytes.NewBufferString("test")},
			test.filepath)
		if test.fail == (err == nil) {
			// t.Errorf("Should test fail? %t. Got: %v", test.fail, err)
		}
	}
}

func TestAuthURL(t *testing.T) {
	tmpG := GoogleDrive{}
	if tmpG.AuthURL() != "" {
		t.Error("Empty config should return empty string")
	}

	if g.AuthURL() == "" {
		t.Error("Configured client should return a value")
	}
}

func TestNewGoogleDrive(t *testing.T) {
	if _, err := NewGoogleDrive(""); err == nil {
		t.Error("Function shoud fail if no filepath is configured")
	}
}

func TestValidate(t *testing.T) {
	if _, _, err := g.Validate(""); err == nil {
		t.Error("Validation should fail with empty token")
	}
}
