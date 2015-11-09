package googledrive

import (
	"fmt"
	"strings"

	"google.golang.org/api/drive/v2"
)

func getFileByPath(sv *drive.Service, path string) (*drive.File, error) {
	paths := strings.Split(sanitize(path), "/")
	title := paths[len(paths)-1]

	parent, err := getParent(sv, path)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("'%s' in parents and title = '%s' and trashed = false",
		parent, title)

	list, err := sv.Files.List().Q(query).Do()
	if err != nil {
		return nil, err
	}

	if len(list.Items) == 0 {
		err := fmt.Errorf("Could not find file %s", title)
		return nil, err
	}

	return list.Items[0], nil
}

func getParent(sv *drive.Service, path string) (*drive.ParentReference, error) {
	paths := strings.Split(sanitize(path), "/")
	if len(paths) == 1 {
		return &drive.ParentReference{Id: "root"}, nil
	}

	parentPath := strings.Join(append(paths[:len(paths)-1]), "/")
	parent, err := getFileByPath(sv, parentPath)
	if err != nil {
		return nil, err
	}

	return &drive.ParentReference{Id: parent.Id}, nil
}

func sanitize(in string) string {
	return strings.TrimPrefix(strings.TrimSuffix(in, "/"), "/")
}
