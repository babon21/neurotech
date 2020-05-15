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

type Publication struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Year  int32         `json:"year" bson:"year"`
	Title string        `json:"title" bson:"title"`
}

type YearPublication struct {
	Year  int32    `json:"year" bson:"_id"`
	Title []string `json:"titles" bson:"titles"`
}

type PublicationHandler struct {
	Collection *mgo.Collection
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
	request.DeleteById(w, r, h.Collection)
}

func (h *PublicationHandler) UpdatePublication(w http.ResponseWriter, r *http.Request) {
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

	jsonNews, err := json.Marshal(publication)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonNews)
	fmt.Println("Edit publication success!")
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
	request.GetOne(w, h.Collection, id, &publication)
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
