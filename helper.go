package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v2"
)

func getClient(r *http.Request) (*drive.Service, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, errors.New("No cookie present in Request")
	}

	token := &oauth2.Token{
		AccessToken: cookie.Value,
	}

	return drive.New(config.Client(context.Background(), token))
}

func getDirectoryList(srv *drive.Service, path string) ([]*drive.File, error) {
	folder, err := getFileByPath(srv, path)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("'%s' in parents and trashed = false", folder.Id)
	list, err := srv.Files.List().Q(query).Do()
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

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
