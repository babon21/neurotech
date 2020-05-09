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

type News struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Title   string        `json:"title" bson:"title"`
	Content string        `json:"content" bson:"content"`
}

type NewsHandler struct {
	Collection *mgo.Collection
}

func (h *NewsHandler) CreateNews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create news request!")

	var news News
	request.Decode(w, r.Body, &news)
	news.ID = bson.NewObjectId()
	request.CreateOne(w, h.Collection, &news)
	
	fmt.Println("Create news success!")
}

func (h *NewsHandler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete news request!")
	request.DeleteById(w, r, h.Collection)
}

func (h *NewsHandler) UpdateNews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit news request!")

	// TODO проверка на bson id

	var news News
	request.Decode(w, r.Body, &news)

	err := h.Collection.Update(bson.M{"_id": news.ID}, news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonNews, err := json.Marshal(news)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonNews)
	fmt.Println("Edit news success!")
}

func (h *NewsHandler) GetNewsList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get news list request!")
	news := []*News{}
	rangeParam := r.URL.Query().Get("range")
	// need check to rangeParam
	request.GetList(w, h.Collection, rangeParam, &news)
}

func (h *NewsHandler) GetOneNews(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("Get one news request! id: ", id)
	news := &News{}
	request.GetOne(w, h.Collection, id, &news)
}

func InitNewsCollection(database *mgo.Database) *mgo.Collection {
	// если коллекции не будет, то она создасться автоматически
	collection := database.C("news")

	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&News{
			bson.NewObjectId(),
			"mongodb",
			"Рассказать про монгу",
		})
		collection.Insert(&News{
			bson.NewObjectId(),
			"redis",
			"Рассказать про redis",
		})
	}

	return collection
}
