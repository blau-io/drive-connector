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
}

func TestGetParent(t *testing.T) {
	service, _ := drive.New(http.DefaultClient)

	if _, err := getParent(nil, ""); err == nil {
		t.Error("Should fail with missing configuration")
	}

	if parent, _ := getParent(service, ""); parent.Id != "root" {
		t.Errorf("Wanted parent with id = root, got %s", parent.Id)
	}
}
