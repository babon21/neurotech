package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/babon21/neurotech/backend/request"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Discipline struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name" bson:"name"`
}

type DisciplineHandler struct {
	Collection *mgo.Collection
}

func (h *DisciplineHandler) CreateDiscipline(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create discipline request!")

	var discipline Discipline
	request.Decode(w, r.Body, &discipline)
	discipline.ID = bson.NewObjectId()
	request.CreateOne(w, h.Collection, &discipline)

	fmt.Println("Create discipline success!")
}

func (h *DisciplineHandler) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete discipline request!")
	request.DeleteById(w, r, h.Collection)
}

func (h *DisciplineHandler) UpdateDiscipline(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit discipline request!")

	// TODO проверка на bson id

	var discipline Discipline
	request.Decode(w, r.Body, &discipline)

	err := h.Collection.Update(bson.M{"_id": discipline.ID}, discipline)
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

func InitDisciplineCollection(database *mgo.Database) *mgo.Collection {
	// если коллекции не будет, то она создасться автоматически
	collection := database.C("discipline")

	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&Discipline{
			bson.NewObjectId(),
			"mongodb discipline",
		})
		collection.Insert(&Discipline{
			bson.NewObjectId(),
			"redis discipline",
		})
	}

	return collection
}

// const DisciplinePath = "disciplines/"

// type DisciplineHandlerOld struct {
// 	path string
// }

// type DisciplinePost struct {
// 	Name string
// }

// type DisciplinePut struct {
// 	OldName string `json:"old_name"`
// 	NewName string `json:"new_name"`
// }

// type DisciplineDelete struct {
// 	Name string
// }

// func (h *DisciplineHandlerOld) CreateDiscipline(w http.ResponseWriter, r *http.Request) {
// 	var post DisciplinePost
// 	err := json.NewDecoder(r.Body).Decode(&post)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	result := utils.CreateFolder(h.path + post.Name)
// 	if result {
// 		fmt.Println("Discipline created!")
// 	} else {
// 		fmt.Println("Failed create discipline!")
// 	}
// }

// func (h *DisciplineHandlerOld) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
// 	var delete DisciplineDelete
// 	err := json.NewDecoder(r.Body).Decode(&delete)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	err2 := os.Remove(h.path + delete.Name)
// 	if err != nil {
// 		fmt.Println(err2)
// 		return
// 	}
// 	fmt.Println("Discipline deleted!")
// }

// func (h *DisciplineHandlerOld) RenameDiscipline(w http.ResponseWriter, r *http.Request) {
// 	var put DisciplinePut
// 	err := json.NewDecoder(r.Body).Decode(&put)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Println(put)

// 	err2 := os.Rename(h.path+put.OldName, h.path+put.NewName)
// 	if err != nil {
// 		fmt.Println(err2)
// 		return
// 	}
// 	fmt.Println("Discipline renamed!")
// }

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
