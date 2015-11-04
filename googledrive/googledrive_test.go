package googledrive

import (
	//	"google.golang.org/api/drive/v2"
	"testing"
)

func TestAdd(t *testing.T) {
	if configured {
		t.SkipNow()
	}

	if Add("", nil, "") == nil {
		t.Error("Should fail with missing configuration")
	}
}

func TestAuthURL(t *testing.T) {
	if configured {
		t.SkipNow()
	}

	if AuthURL() != "" {
		t.Errorf("Should fail with missing configuration")
	}
}

func TestConfig(t *testing.T) {
	if Config("") == nil {
		t.Error("Config should throw error when an empty filepath is used")
	}
}

func TestValidate(t *testing.T) {
	if configured {
		t.SkipNow()
	}

	if _, _, err := Validate("foo"); err == nil {
		t.Error("Should fail with missing configuration")
	}
}
