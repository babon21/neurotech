package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/babon21/neurotech/backend/request"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type StudentWork struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Student     string
	Year        int32  `json:"year" bson:"year"`
	Description string `json:"description" bson:"description"`
}

type StudentWorkHandler struct {
	Collection *mgo.Collection
}

func (h *StudentWorkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateStudentWork(w, r)
	case http.MethodDelete:
		h.DeleteStudentWork(w, r)
	case http.MethodPut:
		h.EditStudentWork(w, r)
	case http.MethodGet:
		h.GetStudentWorkList(w, r)
	}
}

func (h *StudentWorkHandler) CreateStudentWork(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create student work request!")

	var studentWork StudentWork
	err := json.NewDecoder(r.Body).Decode(&studentWork)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	studentWork.ID = bson.NewObjectId()
	err2 := h.Collection.Insert(studentWork)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Create student work success!")
}

func (h *StudentWorkHandler) DeleteStudentWork(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete student work request!")

	// TODO проверка на bson id
	var delete request.DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&delete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := h.Collection.Remove(bson.M{"_id": delete.ID})
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Delete student work success!")
}

func (h *StudentWorkHandler) EditStudentWork(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Edit student work success!")
}

func (h *StudentWorkHandler) GetStudentWorkList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get student work list request!")

	studentWork := []*StudentWork{}
	// bson.M{} - это типа условия для поиска
	err := h.Collection.Find(bson.M{}).All(&studentWork)
	if err != nil {
		panic(err)
	}

	jsonStudentWorks, err := json.Marshal(studentWork)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Get student work list success!")

	ResponseWithJSON(w, jsonStudentWorks)
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
		})
		collection.Insert(&StudentWork{
			bson.NewObjectId(),
			"Some Batman",
			2014,
			"Диплом про redis",
		})
	}

	return collection
}
