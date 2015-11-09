package googledrive

import (
	"net/http"
	"testing"

	"google.golang.org/api/drive/v2"
)

func TestGetParent(t *testing.T) {
	service, _ := drive.New(http.DefaultClient)

	parent, err := getParent(service, "")
	if err != nil {
		t.Error(err.Error())
	}

	if parent.Id != "root" {
		t.Errorf("Wanted parent id = 'root', got %s", parent.Id)
	}
}

func TestSanitize(t *testing.T) {
	var sanitzeTable = []struct {
		in  string
		out string
	}{
		{"", ""},
		{"/", ""},
		{"foo", "foo"},
		{"/foo", "foo"},
		{"/foo/", "foo"},
		{"/foo/bar", "foo/bar"},
		{"/foo/bar/", "foo/bar"},
	}

	for _, test := range sanitzeTable {
		if sanitize(test.in) != test.out {
			t.Errorf("Expected %s, got %s", test.in, sanitize(test.out))
		}
	}
}
