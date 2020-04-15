package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Publication struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Year    int32         `json:"year" bson:"year"`
	Content string        `json:"content" bson:"content"`
}

type PublicationHandler struct {
	Publication *mgo.Collection
}

type PublicationDelete struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

func (h *PublicationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreatePublication(w, r)
	case http.MethodDelete:
		h.DeletePublication(w, r)
	case http.MethodPut:
		h.EditPublication(w, r)
	case http.MethodGet:
		h.GetPublications(w, r)
	}
}

func (h *PublicationHandler) CreatePublication(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create publication request!")

	var publication Publication
	err := json.NewDecoder(r.Body).Decode(&publication)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	publication.ID = bson.NewObjectId()
	err2 := h.Publication.Insert(publication)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Create publication success!")
}

func (h *PublicationHandler) DeletePublication(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete publication request!")

	// TODO проверка на bson id
	var delete PublicationDelete
	err := json.NewDecoder(r.Body).Decode(&delete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := h.Publication.Remove(bson.M{"_id": delete.ID})
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Delete publication success!")
}

func (h *PublicationHandler) EditPublication(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit publication request!")

	// TODO проверка на bson id

	var publication Publication
	err := json.NewDecoder(r.Body).Decode(&publication)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := h.Publication.Update(bson.M{"_id": publication.ID}, publication)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Edit publication success!")
}

func (h *PublicationHandler) GetPublications(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get publications request!")

	publication := []*Publication{}
	// bson.M{} - это типа условия для поиска
	err := h.Publication.Find(bson.M{}).All(&publication)
	if err != nil {
		panic(err)
	}

	jsonNews, err := json.Marshal(publication)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Get publications success!")

	ResponseWithJSON(w, jsonNews)
}
