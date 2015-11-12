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

var emptyG = GoogleDrive{}

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
		//{"", "test", true},
		{"token", "", true},
		//{"token", "test", false},
	}

	for _, test := range addTestTable {
		err := g.Add(test.code, test.filepath,
			mockRC{bytes.NewBufferString("test")})
		if test.fail == (err == nil) {
			t.Errorf("Error expected: %t. Got: %v", test.fail, err)
		}
	}
}

func TestAuthURL(t *testing.T) {
	if emptyG.AuthURL() != "" {
		t.Error("Empty config should not return anything")
	}

	if g.AuthURL() == "" {
		t.Error("Configured client should return a value")
	}
}

func TestBrowse(t *testing.T) {
	if list, _ := emptyG.Browse("", ""); list != nil {
		t.Error("Empty config should not return anything")
	}
}

func TestRemove(t *testing.T) {
	if emptyG.Remove("", "") == nil {
		t.Error("You should not be able to remove the root folder")
	}

	if emptyG.Remove("", "test") != nil {
		t.Error("Empty config should just fail silently")
	}
}

func TestNewGoogleDrive(t *testing.T) {
	if _, err := NewGoogleDrive(""); err == nil {
		t.Error("Function shoud fail if no filepath is configured")
	}
}

func TestPublish(t *testing.T) {
	if _, err := g.Publish("", "/"); err == nil {
		t.Error("Shouldn't be able to publish root folder")
	}

	if _, err := emptyG.Publish("", "test"); err != nil {
		t.Error("Empty config should fail silently")
	}
}

func TestRead(t *testing.T) {
	if _, err := g.Read("", "/"); err == nil {
		t.Error("Shouldn't be able to read a directory")
	}

	if _, err := emptyG.Read("", "test"); err != nil {
		t.Error("Empty config should fail silently")
	}
}

func TestValidate(t *testing.T) {
	if _, _, err := emptyG.Validate(""); err != nil {
		t.Error("Empty config should fail silently")
	}

	if _, _, err := g.Validate(""); err == nil {
		t.Error("Validation should fail with empty token")
	}

	if _, _, err := g.Validate("invalid"); err == nil {
		t.Error("an invalid code should fail")
	}
}
