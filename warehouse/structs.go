package warehouse

import (
	"time"
)

// AuthURLJSON is a json struct which returns the link to an oauth2 provider
type AuthURLJSON struct {
	URL string `json:"url"`
}

// BrowseJSON is a json struct which return a list of files for a given folder
type BrowseJSON struct {
	FileList []string `json:"file_list"`
}

// PublishJSON is a json struct which returns the public link to a file
type PublishJSON struct {
	URL string `json:"url"`
}

// ValidateJSON is a json struct which returns the auth token and it's
// expiry date
type ValidateJSON struct {
	Token  string    `json:"access_token"`
	Expiry time.Time `json:"expiry,omitempty"`
}
