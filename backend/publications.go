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

type Publication struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Year    int32         `json:"year" bson:"year"`
	Content string        `json:"content" bson:"content"`
}

type PublicationHandler struct {
	Collection *mgo.Collection
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
	err2 := h.Collection.Insert(publication)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Create publication success!")
}

func (h *PublicationHandler) DeletePublication(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete publication request!")

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

	err2 := h.Collection.Update(bson.M{"_id": publication.ID}, publication)
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
	err := h.Collection.Find(bson.M{}).All(&publication)
	if err != nil {
		panic(err)
	}

	jsonPublications, err := json.Marshal(publication)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Get publications success!")

	ResponseWithJSON(w, jsonPublications)
}

func InitPublicationsCollection(database *mgo.Database) *mgo.Collection {
	// если коллекции не будет, то она создасться автоматически
	collection := database.C("publications")

	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&Publication{
			bson.NewObjectId(),
			2018,
			"Публикация про монгу",
		})
		collection.Insert(&Publication{
			bson.NewObjectId(),
			2019,
			"Публикация про redis",
		})
	}

	return collection
}