package main

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/drive/v2"
)

func getFileByPath(srv *drive.Service, path string) (*drive.File, error) {
	paths := strings.Split(strings.TrimPrefix(path, "/"), "/")
	parent := "root"
	var err error

	for i := 0; i < len(paths)-1; i++ {
		query := fmt.Sprintf("'%s' in parents and title = '%s' and mimeType = "+
			"'application/vnd.google-apps.folder' and trashed = false", parent,
			paths[i])
		list, err := srv.Files.List().Q(query).Do()
		if err != nil {
			return nil, err
		}

		if len(list.Items) == 0 {
			err = errors.New(fmt.Sprintf("Could not find file %s", paths[i]))
			return nil, err
		}

		parent = list.Items[0].Id
	}

	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("'%s' in parents and title = '%s' and trashed = false",
		parent, paths[len(paths)-1])
	list, err := srv.Files.List().Q(query).Do()
	if err != nil {
		return nil, err
	}

	if len(list.Items) == 0 {
		return nil, errors.New(fmt.Sprintf("Could not find file %s",
			paths[len(paths)-1]))
	}

	return list.Items[0], nil
}
