package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var uploadFormTmpl = []byte(`
	<html>
	<body>
	<form action="/study" method="post" enctype="multipart/form-data">
	Image: <input type="file" name="filename">
	<input type="submit" value="Upload">
	</form>
	<form action="/study/delete" method="post">
	text: <input type="text" name="filename">
	<input type="submit" value="submit">
	</form>
	</body>
	</html>
	`)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(uploadFormTmpl)
	})

	http.HandleFunc("/study", studyHandler)
	http.HandleFunc("/study/delete", deleteHandler)
	http.ListenAndServe(":8080", nil)
}

func studyHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		err := uploadFile(r, "./study", 5*1024*1025)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func uploadFile(r *http.Request, path string, size int64) error {
	r.ParseMultipartForm(size)

	file, handler, err := r.FormFile("filename")
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := os.OpenFile(path+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	err := os.Remove("./study/" + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
}
