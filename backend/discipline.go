package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/babon21/neurotech/backend/utils"
	"github.com/rs/zerolog/log"
)

const DisciplinePath = "disciplines/"

type DisciplineHandler struct {
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

func (h *DisciplineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *DisciplineHandler) CreateDiscipline(w http.ResponseWriter, r *http.Request) {
	var post DisciplinePost
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := utils.CreateFolder(h.path + post.Name)
	if result {
		fmt.Println("Discipline created!")
	} else {
		fmt.Println("Failed create discipline!")
	}
}

func (h *DisciplineHandler) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
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

func (h *DisciplineHandler) RenameDiscipline(w http.ResponseWriter, r *http.Request) {
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

func (h *DisciplineHandler) GetDisciplineList(w http.ResponseWriter, r *http.Request) {
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
