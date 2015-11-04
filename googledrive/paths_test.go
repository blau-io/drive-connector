package googledrive

import (
	"net/http"
	"testing"

	"google.golang.org/api/drive/v2"
)

func TestGetFileByPath(t *testing.T) {
	if _, err := getFileByPath(nil, ""); err == nil {
		t.Error("Should fail with missing configuration")
	}

	var fileTable = []struct {
		filepath string
		title    string
	}{
	//		{"index.html", "index.html"},
	}

	//service, _ := drive.New(http.DefaultClient)

	for _, test := range fileTable {
		file, err := getFileByPath(nil, test.filepath)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		if file.Title != test.title {
			t.Errorf("Expected title %s, got %s", test.title, file.Title)
		}
	}

}

func TestGetParent(t *testing.T) {
	if _, err := getParent(nil, ""); err == nil {
		t.Error("Should fail with missing configuration")
	}

	var parentTable = []struct {
		filepath string
		parent   string
	}{
		{"", "root"},
		//	{"/foo/bar/index.html", "bar"},
	}

	service, _ := drive.New(http.DefaultClient)

	for _, test := range parentTable {
		parent, err := getParent(service, test.filepath)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		if parent.Id != test.parent {
			t.Errorf("Wanted parent id = %s, got %s", test.filepath, parent.Id)
		}
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
