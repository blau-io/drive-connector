package warehouse

import (
	"testing"
)

var c, _ = NewClient("")

func TestAdd(t *testing.T) {
	if c.Add("", "", nil) == nil {
		t.Error("Empty filepath should throw error")
	}
}

func TestAuthURL(t *testing.T) {
	if _, err := c.AuthURL(""); err == nil {
		t.Error("Empty provider should throw error")
	}

	if _, err := c.AuthURL("google"); err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestBrowse(t *testing.T) {
	if _, err := c.Browse("", "test"); err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestNewClient(t *testing.T) {
	if clt, _ := NewClient(""); clt == nil {
		t.Error("NewClient should never return nil client")
	}
}

func TestPublish(t *testing.T) {
	if _, err := c.Publish("", ""); err == nil {
		t.Error("Shoudn't be able to publish root folder")
	}
}

func TestRead(t *testing.T) {
	if _, err := c.Read("", ""); err == nil {
		t.Error("Empty filepath should throw error")
	}

	if _, err := c.Read("", "test"); err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestRemove(t *testing.T) {
	if c.Remove("", "") == nil {
		t.Error("Should not be able to remove root folder")
	}
}

func TestValidate(t *testing.T) {
	if _, _, err := c.Validate("", ""); err == nil {
		t.Error("Empty state should throw error")
	}

	if _, _, err := c.Validate("google", ""); err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}
