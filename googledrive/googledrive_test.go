package googledrive

import (
	//	"google.golang.org/api/drive/v2"
	"testing"
)

func TestAdd(t *testing.T) {
	//if Add("", nil, "") == nil {
	//	t.Error("Bad config should throw error")
	//}
}

func TestConfig(t *testing.T) {
	if Config("") == nil {
		t.Error("Config should throw error when an empty filepath is used")
	}
}
