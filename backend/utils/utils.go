package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

var singleFileUploadForm = []byte(`
	<html>
	<body>
	<form action="/single/upload" method="post" enctype="multipart/form-data">
	File: <input type="file" name="filename">
	<input type="submit" value="Upload">
	</form>
	</body>
	</html>
	`)

var multipleFileUploadForm = []byte(`
	<html>
	<body>
	<form action="/multiple/upload" enctype="multipart/form-data" method="post">
	<input type="file" name="multiplefiles" id="multiplefiles" multiple>
	<input type="submit" name="submit" >
	</form>
	</body>
	</html>
`)

func UploadMultipleFilesFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(multipleFileUploadForm))
}

func UploadSingleFileFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(singleFileUploadForm))
}

func UploadFile(r *http.Request, path string, size int64) error {
	fmt.Println("upload single file start")
	r.ParseMultipartForm(size)

	file, handler, err := r.FormFile("filename")
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := os.OpenFile(path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func CreateFolder(folderName string) bool {
	_, err := os.Stat(folderName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderName, 0755)
		if errDir != nil {
			log.Err(errDir).Msg("mkdir err")
			return false
		}
		return true
	}
	return false
}
