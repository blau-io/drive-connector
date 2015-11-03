package googledrive

import (
	"testing"
)

func TestConfig(t *testing.T) {
	if Config("") == nil {
		t.Error("Config should throw error when an empty filepath is used")
	}
}
