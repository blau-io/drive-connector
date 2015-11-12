package googledrive

import (
	"golang.org/x/oauth2"
	"testing"
)

var mockConfig = &oauth2.Config{
	ClientID:     "CLIENT_ID",
	ClientSecret: "CLIENT_SECRET",
	RedirectURL:  "REDIRECT_URL",
	Scopes:       []string{"scope1", "scope2"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "http://example.com/auth",
		TokenURL: "http://example.com/token",
	},
}

func TestGetClient(t *testing.T) {
	var clientTable = []struct {
		config *oauth2.Config
		fail   bool
	}{
		{nil, true},
		{mockConfig, false},
	}

	for _, test := range clientTable {
		_, err := getClient(test.config, "test")
		if test.fail == (err == nil) {
			t.Errorf("Error: expected %t, got %v", test.fail, err)
		}
	}
}

func TestIsRoot(t *testing.T) {
	var isRootTable = []struct {
		filepath string
		result   bool
	}{
		{"", true},
		{"/", true},
		{"foo", false},
	}

	for _, test := range isRootTable {
		result := isRoot(test.filepath)

		if test.result != result {
			t.Errorf("Error: expected %t, got %t", test.result, result)
		}
	}
}
