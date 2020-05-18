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

type StudentWork struct {
	ID      bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Student string        `json:"student" bson:"student"`
	Year    int32         `json:"year" bson:"year"`
	Title   string        `json:"title" bson:"title"`
	Type    string        `json:"type,omitempty" bson:"type"`
}

type StudentWorkHandler struct {
	Collection *mgo.Collection
}

func (h *StudentWorkHandler) CreateStudentWork(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create student work request!")

	var studentWork StudentWork
	request.Decode(w, r.Body, &studentWork)
	studentWork.ID = bson.NewObjectId()
	request.CreateOne(w, h.Collection, &studentWork)

	fmt.Println("Create student work success!")
}

func (h *StudentWorkHandler) DeleteStudentWork(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete student work request!")

	request.DeleteById(w, r, h.Collection)

	fmt.Println("Delete student work success!")
}

func (h *StudentWorkHandler) UpdateStudentWork(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit student work request!")

	// TODO проверка на bson id

	var studentWork StudentWork
	err := json.NewDecoder(r.Body).Decode(&studentWork)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := h.Collection.Update(bson.M{"_id": studentWork.ID}, studentWork)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonWork, err := json.Marshal(studentWork)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonWork)
	fmt.Println("Edit student work success!")
}

func (h *StudentWorkHandler) GetStudentWorkList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get student work list request!")
	typeWork := r.URL.Query().Get("type_group")

	if typeWork == "" {
		list := []*StudentWork{}
		rangeParam := r.URL.Query().Get("range")
		// need check to rangeParam
		request.GetList(w, h.Collection, rangeParam, &list)
		return
	}

	fmt.Println("Get publications list with year group request from site!")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type": typeWork,
			},
		},
		{
			"$sort": bson.M {
				"year": -1,
			},
		},
	}

	var result []StudentWork

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
	fmt.Println("Get student work list success!")
}

func (h *StudentWorkHandler) GetOneStudentWork(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("Get one student work request! id: ", id)
	work := &StudentWork{}
	request.GetOne(w, h.Collection, id, &work)
}

func InitStudentWorksCollection(database *mgo.Database) *mgo.Collection {
	// если коллекции не будет, то она создасться автоматически
	collection := database.C("student_works")

	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&StudentWork{
			bson.NewObjectId(),
			"Иванов Иван",
			2017,
			"Диплом про монгу",
			"bachelor_work",
		})
		collection.Insert(&StudentWork{
			bson.NewObjectId(),
			"Some Batman",
			2014,
			"Диплом про redis",
			"master_dissertation",
		})
	}

	return collection
}
