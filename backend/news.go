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

type News struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Title   string        `json:"title" bson:"title"`
	Content string        `json:"content" bson:"content"`
}

type NewsHandler struct {
	Collection *mgo.Collection
}

func (h *NewsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateNews(w, r)
	case http.MethodDelete:
		h.DeleteNews(w, r)
	case http.MethodPut:
		h.PutNews(w, r)
	case http.MethodGet:
		h.GetNewsList(w, r)
	}
}

func (h *NewsHandler) CreateNews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create news request!")

	var news News
	err := json.NewDecoder(r.Body).Decode(&news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	news.ID = bson.NewObjectId()
	err2 := h.Collection.Insert(news)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Create news success!")
}

func (h *NewsHandler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete news request!")

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

	fmt.Println("Delete news success!")
}

func (h *NewsHandler) PutNews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit news request!")

	// TODO проверка на bson id

	var news News
	err := json.NewDecoder(r.Body).Decode(&news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := h.Collection.Update(bson.M{"_id": news.ID}, news)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Edit news success!")
}

func (h *NewsHandler) GetNewsList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get news list request!")

	news := []*News{}
	// bson.M{} - это типа условия для поиска
	err := h.Collection.Find(bson.M{}).All(&news)
	if err != nil {
		panic(err)
	}

	jsonNews, err := json.Marshal(news)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Get news list success!")

	ResponseWithJSON(w, jsonNews)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(json)
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
