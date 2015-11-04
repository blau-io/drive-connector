package googledrive

import (
	"errors"
	"strings"

	"google.golang.org/api/drive/v2"
)

func getFileByPath(srv *drive.Service, path string) (*drive.File, error) {
	return nil, errors.New("TODO")
}

func getParent(sv *drive.Service, path string) (*drive.ParentReference, error) {
	if sv == nil {
		return nil, errors.New("No drive.Service pointer passed")
	}

	paths := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(paths) <= 1 {
		return &drive.ParentReference{Id: "root"}, nil
	}

	parentPath := strings.Join(append(paths[:len(paths)-1]), "/")
	parent, err := getFileByPath(sv, parentPath)
	if err != nil {
		return nil, err
	}

	return &drive.ParentReference{Id: parent.Id}, nil
}
