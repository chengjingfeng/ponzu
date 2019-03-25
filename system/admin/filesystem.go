package admin

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/ponzu-cms/ponzu/system/db"
	"github.com/ponzu-cms/ponzu/system/item"
)

func deleteUploadFromDisk(target string) error {
	// get data on file
	data, err := db.Upload(target)
	if err != nil {
		return err
	}

	// unmarshal data
	upload := item.FileUpload{}
	if err = json.Unmarshal(data, &upload); err != nil {
		return err
	}

	// use path to delete the physical file from disk
	delPath := strings.Replace(upload.Path, "/api/", "./", 1)
	err = os.Remove(delPath)
	if err != nil {
		return err
	}

	return nil
}

func restrict(dir http.Dir) justFilesFilesystem {
	return justFilesFilesystem{dir}
}

// the code below removes the open directory listing when accessing a URL which
// normally would point to a directory. code from golang-nuts mailing list:
// https://groups.google.com/d/msg/golang-nuts/bStLPdIVM6w/hidTJgDZpHcJ
// credit: Brad Fitzpatrick (c) 2012

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
