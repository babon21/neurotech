package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/babon21/neurotech/backend/request"
	"github.com/babon21/neurotech/backend/utils"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Discipline struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Files []File        `json:"files" bson:"files"`
}

type File struct {
	Url  string `json:"url" bson:"url"`
	Name string `json:"name" bson:"name"`
}

type DisciplineHandler struct {
	Collection *mgo.Collection
	Path       string
}

func (h *DisciplineHandler) CreateDiscipline(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create discipline request!")

	var discipline Discipline
	request.Decode(w, r.Body, &discipline)
	discipline.ID = bson.NewObjectId()
	request.CreateOne(w, h.Collection, &discipline)

	result := utils.CreateFolder(h.Path + discipline.Name)
	if result {
		fmt.Println("Discipline folder created!")
	} else {
		fmt.Println("Failed create folder discipline!")
	}

	fmt.Println("Create discipline success!")
}

func (h *DisciplineHandler) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete discipline request!")
	id := chi.URLParam(r, "id")

	var discipline Discipline
	err := h.Collection.FindId(bson.ObjectIdHex(id)).One(&discipline)

	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(h.Path + discipline.Name)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	request.DeleteById(w, r, h.Collection)

	fmt.Println("Success delete discipline request!")
}

func (h *DisciplineHandler) UpdateDiscipline(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit discipline request!")
	id := chi.URLParam(r, "id")

	// TODO проверка на bson id

	var discipline Discipline

	err := h.Collection.FindId(bson.ObjectIdHex(id)).One(&discipline)
	if err != nil {
		panic(err)
	}

	prevName := discipline.Name
	prevFiles := discipline.Files

	request.Decode(w, r.Body, &discipline)

	removedFiles := getRemovedFiles(prevFiles, discipline.Files)
	if removedFiles != nil {
		removeFiles(h.Path + discipline.Name + "/", removedFiles)
	}

	if prevName != discipline.Name {
		err = os.Rename(h.Path+prevName, h.Path+discipline.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = h.Collection.Update(bson.M{"_id": discipline.ID}, discipline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonNews, err := json.Marshal(discipline)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonNews)
	fmt.Println("Edit discipline success!")
}

func (h *DisciplineHandler) GetDisciplineList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get news list request!")
	disciplines := []*Discipline{}
	rangeParam := r.URL.Query().Get("range")
	// need check to rangeParam
	request.GetList(w, h.Collection, rangeParam, &disciplines)
}

func (h *DisciplineHandler) GetOneDiscipline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("Get one discipline request! id: ", id)
	discipline := &Discipline{}
	request.GetOne(w, h.Collection, id, &discipline)
}

func (h *DisciplineHandler) UploadDisciplineFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Upload discipline files request!")
	id := chi.URLParam(r, "id")
	fmt.Println("Upload files to discipline request! id: ", id)

	// TODO проверка на bson id

	var discipline Discipline

	err := h.Collection.FindId(bson.ObjectIdHex(id)).One(&discipline)
	if err != nil {
		panic(err)
	}

	// в папку этой дисциплины
	filenames, err := h.UploadMultipleFiles(w, r, 5*1024*1025*100, discipline.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := formDicsiplineFiles(filenames, discipline.Name)
	discipline.Files = append(discipline.Files, files...)

	err = h.Collection.Update(bson.M{"_id": discipline.ID}, discipline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Upload discipline files success!")
}

func formDicsiplineFiles(filenames []string, name string) []File {
	var files []File
	for _, f := range filenames {
		file := File{Url: name + "/" + f, Name: f}
		files = append(files, file)
	}

	return files
}

func (h *DisciplineHandler) UploadMultipleFiles(w http.ResponseWriter, r *http.Request, size int64, name string) ([]string, error) {
	r.ParseMultipartForm(size)
	path := h.Path + name

	files := r.MultipartForm.File["multiplefiles"] // grab the filenames

	var filenames []string
	// loop through the files one by one
	for _, f := range files {
		file, err := f.Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return nil, err
		}

		log.Info().Msg("add study material to discipline name: " + f.Filename)
		out, err := os.Create(path + "/" + f.Filename)

		defer out.Close()
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return nil, err
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			fmt.Fprintln(w, err)
			return nil, err
		}

		filenames = append(filenames, f.Filename)
		log.Info().Msg("Success uploaded file:" + f.Filename)
	}

	return filenames, nil
}

func InitDisciplineCollection(database *mgo.Database) *mgo.Collection {
	// если коллекции не будет, то она создасться автоматически
	collection := database.C("discipline")

	var files1 []File
	files1 = append(files1, File{Url: "/disc1/book1", Name: "book1.pdf"})
	files1 = append(files1, File{Url: "/disc1/book2", Name: "book2.pdf"})

	var files2 []File
	files2 = append(files2, File{Url: "/disc2/book3", Name: "book3.pdf"})
	files2 = append(files2, File{Url: "/disc2/book4", Name: "book4.pdf"})

	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&Discipline{
			bson.NewObjectId(),
			"mongodb discipline",
			files1,
		})
		collection.Insert(&Discipline{
			bson.NewObjectId(),
			"redis discipline",
			files2,
		})
	}

	return collection
}

func getRemovedFiles(prevFiles []File, currentFiles []File) []File {
	var removedFiles []File

	for _, f := range prevFiles {
		if !fileExists(currentFiles, f) {
			removedFiles = append(removedFiles, f)
		}
	}

	if len(removedFiles) == 0 {
		return nil
	}

	return removedFiles
}

func fileExists(files []File, file File) bool {
	for _, f := range files {
		if file == f {
			return true
		}
	}

	return false
}

func removeFiles(prefix string, files []File) {
	for _, f := range files {
		err := os.Remove(prefix + f.Name)
		if err != nil {
			panic(err)
		}
	}
}

// func (h *DisciplineHandlerOld) GetDisciplineList(w http.ResponseWriter, r *http.Request) {
// 	files, err := ioutil.ReadDir(h.path)
// 	if err != nil {
// 		log.Err(err).Msg("read dir err")
// 		return
// 	}

// 	fileNames := make([]string, 0, len(files))
// 	for _, file := range files {
// 		fmt.Println(file.Name())
// 		fileNames = append(fileNames, file.Name())
// 	}

// 	fmt.Println(fileNames)
// 	fileNamesJSON, err := json.Marshal(fileNames)

// 	if err != nil {
// 		log.Err(err).Msg("json marshall err")
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(fileNamesJSON)
// }
