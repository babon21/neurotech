package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
	disciplineHandler := &Handler{path: "disciplines/"}
	http.Handle("/disciplines", disciplineHandler)
	http.ListenAndServe(":8080", nil)
}

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

type Handler struct {
	path string
}

type DisciplinePost struct {
	Name string
}

type DisciplinePut struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type DisciplineDelete struct {
	Name string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("Post request!")
		h.CreateDiscipline(w, r)
	case http.MethodDelete:
		fmt.Println("Delete request!")
		h.DeleteDiscipline(w, r)
	case http.MethodPut:
		fmt.Println("Put request!")
		h.RenameDiscipline(w, r)
	case http.MethodGet:
		fmt.Println("Get request!")
		h.GetDisciplineList(w, r)
	}
}

func (h *Handler) CreateDiscipline(w http.ResponseWriter, r *http.Request) {
	var post DisciplinePost
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := CreateFolder(h.path + post.Name)
	if result {
		fmt.Println("Discipline created!")
	} else {
		fmt.Println("Failed create discipline!")
	}
}

func (h *Handler) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
	var delete DisciplineDelete
	err := json.NewDecoder(r.Body).Decode(&delete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := os.Remove(h.path + delete.Name)
	if err != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("Discipline deleted!")
}

func (h *Handler) RenameDiscipline(w http.ResponseWriter, r *http.Request) {
	var put DisciplinePut
	err := json.NewDecoder(r.Body).Decode(&put)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(put)

	err2 := os.Rename(h.path+put.OldName, h.path+put.NewName)
	if err != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("Discipline renamed!")
}

func (h *Handler) GetDisciplineList(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(h.path)
	if err != nil {
		log.Err(err).Msg("read dir err")
		return
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fmt.Println(file.Name())
		fileNames = append(fileNames, file.Name())
	}

	fmt.Println(fileNames)
	fileNamesJSON, err := json.Marshal(fileNames)

	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(fileNamesJSON)
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
