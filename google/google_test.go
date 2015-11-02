package google

import (
	"testing"
)

func TestAuthUrl(t *testing.T) {
}

func TestConfig(t *testing.T) {
	if Config("") == nil {
		t.Error("Config should throw error when an empty filepath is used")
	}
}
