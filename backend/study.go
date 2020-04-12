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

type StudyMaterialHandler struct {
	path string
}

type StudyMaterialPatch struct {
	OldName        string `json:"old_name"`
	NewName        string `json:"new_name"`
	DisciplineName string `json:"discipline_name"`
}

type StudyMaterialDelete struct {
	Name           string
	DisciplineName string `json:"discipline_name"`
}

func (h *StudyMaterialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		log.Info().Msg("Post request study material")
		h.CreateStudyMaterial(w, r)
	case http.MethodDelete:
		log.Info().Msg("Delete request study material")
		h.DeleteStudyMaterial(w, r)
	case http.MethodPatch:
		log.Info().Msg("Patch request study material")
		h.RenameStudyMaterial(w, r)
	case http.MethodGet:
		log.Info().Msg("Get request study material")
		h.GetStudyMaterials(w, r)
	}
}

func (h *StudyMaterialHandler) CreateStudyMaterial(w http.ResponseWriter, r *http.Request) {
	err := h.UploadMultipleFiles(w, r, 5*1024*1025*100)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *StudyMaterialHandler) DeleteStudyMaterial(w http.ResponseWriter, r *http.Request) {
	var delete StudyMaterialDelete
	err := json.NewDecoder(r.Body).Decode(&delete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info().Msg("delete study material from discipline name: " + delete.DisciplineName)
	err2 := os.Remove(h.path + delete.DisciplineName + "/" + delete.Name)
	if err != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("Study material deleted!")
}

func (h *StudyMaterialHandler) RenameStudyMaterial(w http.ResponseWriter, r *http.Request) {
	var put StudyMaterialPatch
	err := json.NewDecoder(r.Body).Decode(&put)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(put)

	oldPath := h.path + put.DisciplineName + "/" + put.OldName
	newPath := h.path + put.DisciplineName + "/" + put.NewName
	err2 := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("Study material renamed!")
}

func (h *StudyMaterialHandler) GetStudyMaterials(w http.ResponseWriter, r *http.Request) {
	disciplineName := r.URL.Query().Get("discipline_name")
	files, err := ioutil.ReadDir(h.path + disciplineName)
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

func (h *StudyMaterialHandler) UploadMultipleFiles(w http.ResponseWriter, r *http.Request, size int64) error {
	r.ParseMultipartForm(size)
	name := r.PostFormValue("discipline_name")
	path := h.path + name + "/"

	files := r.MultipartForm.File["multiplefiles"] // grab the filenames

	// loop through the files one by one
	for _, f := range files {
		file, err := f.Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return err
		}

		log.Info().Msg("add study material to discipline name: " + f.Filename)
		out, err := os.Create(path + f.Filename)

		defer out.Close()
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return err
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			fmt.Fprintln(w, err)
			return err
		}

		fmt.Fprintf(w, "Files uploaded successfully : ")
		fmt.Fprintf(w, f.Filename+"\n")
	}
	return nil
}
