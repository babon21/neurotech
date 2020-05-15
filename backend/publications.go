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

type Publication struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	Year            int32         `json:"year" bson:"year"`
	Title           string        `json:"title" bson:"title"`
	FilePublication *File         `json:"file,omitempty" bson:"file,omitempty"`
}

type YearPublication struct {
	Year  int32    `json:"year" bson:"_id"`
	Title []string `json:"titles" bson:"titles"`
}

type PublicationHandler struct {
	Collection *mgo.Collection
	Path       string
}

func (h *PublicationHandler) CreatePublication(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create publication request!")

	var publication Publication
	request.Decode(w, r.Body, &publication)
	publication.ID = bson.NewObjectId()
	request.CreateOne(w, h.Collection, &publication)

	fmt.Println("Create publication success!")
}

func (h *PublicationHandler) DeletePublication(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete publication request!")

	id := chi.URLParam(r, "id")

	var publication Publication
	err := h.Collection.FindId(bson.ObjectIdHex(id)).One(&publication)

	if err != nil {
		panic(err)
	}

	if publication.FilePublication != nil {
		os.RemoveAll(h.Path + publication.FilePublication.Name)
	}

	request.DeleteById(w, r, h.Collection)
}

func (h *PublicationHandler) UpdatePublication(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit publication request!")
	id := chi.URLParam(r, "id")
	var prev Publication

	err := h.Collection.FindId(bson.ObjectIdHex(id)).One(&prev)
	if err != nil {
		panic(err)
	}

	// TODO проверка на bson id

	var publication Publication
	err = json.NewDecoder(r.Body).Decode(&publication)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	syncFile(h.Path, prev.FilePublication, publication.FilePublication)

	err = h.Collection.Update(bson.M{"_id": publication.ID}, publication)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonNews, err := json.Marshal(publication)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonNews)
	fmt.Println("Edit publication success!")
}

func syncFile(basePath string, prevFile *File, curFile *File) {
	if prevFile != nil {
		if curFile != nil {
			if *prevFile != *curFile {
				os.Remove(basePath + prevFile.Name)
			}
		} else {
			os.Remove(basePath + prevFile.Name)
		}
	}
}

func (h *PublicationHandler) UploadPublicationFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("Upload files to publication request! id: ", id)

	// TODO проверка на bson id

	var publication Publication

	err := h.Collection.FindId(bson.ObjectIdHex(id)).One(&publication)
	if err != nil {
		panic(err)
	}

	filename, err := UploadFile(w, r, h.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file := &File{
		Url:  h.Path + filename,
		Name: filename,
	}
	publication.FilePublication = file

	err = h.Collection.Update(bson.M{"_id": publication.ID}, publication)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Upload " + publication.Title + " discipline files success!")
}

func UploadFile(w http.ResponseWriter, r *http.Request, path string) (string, error) {
	r.ParseMultipartForm(5 * 1024 * 1025 * 100)
	f, h, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	out, err := os.Create(path + h.Filename)

	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return "", err
	}
	defer out.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return "", err
	}

	_, err = io.Copy(out, f) // file not files[i] !

	log.Info().Msg("Success uploaded file:" + h.Filename)

	return h.Filename, nil
}

func (h *PublicationHandler) GetPublicationsList(w http.ResponseWriter, r *http.Request) {
	isGroup := r.URL.Query().Get("year_group")

	if isGroup == "" {
		fmt.Println("Get publications list request!")
		rangeParam := r.URL.Query().Get("range")
		// need check to rangeParam
		publications := []*Publication{}
		request.GetList(w, h.Collection, rangeParam, &publications)
		return
	}

	fmt.Println("Get publications list with year group request from site!")

	var result []YearPublication
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "$year",
				"titles": bson.M{
					"$push": "$title",
				},
			},
		},
	}

	err := h.Collection.Pipe(pipeline).All(&result)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(result)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	request.ResponseWithJSON(w, json)
}

func (h *PublicationHandler) GetOnePublication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("Get one publications request! id: ", id)
	publication := &Publication{}
	request.GetOne(w, h.Collection, id, publication)
}

func InitPublicationsCollection(database *mgo.Database) *mgo.Collection {
	// если коллекции не будет, то она создасться автоматически
	collection := database.C("publications")

	result := utils.CreateFolder("publications")
	if result {
		fmt.Println("Publications dir created!")
	} else {
		fmt.Println("Failed create publications dir!")
	}

	file1 := &File{Url: "/publications/files/file1", Name: "file1.pdf"}
	file2 := &File{Url: "/publications/files/file2", Name: "file2.pdf"}

	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&Publication{
			bson.NewObjectId(),
			2018,
			"Публикация про монгу",
			file1,
		})
		collection.Insert(&Publication{
			bson.NewObjectId(),
			2019,
			"Публикация про redis",
			file2,
		})
	}

	return collection
}
