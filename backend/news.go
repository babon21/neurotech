package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/babon21/neurotech/backend/utils"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const NewsPath = "news/"

type NewsHandler struct {
	path string
	News *mgo.Collection
}

type NewsPost struct {
	Title string
	Name  string
}

type NewsPut struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Title string
	Name  string
}

type NewsDelete struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

func (h *NewsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("Post request!")
		h.CreateNews(w, r)
	case http.MethodDelete:
		fmt.Println("Delete request!")
		h.DeleteNews(w, r)
	case http.MethodPut:
		fmt.Println("Put request!")
		h.PutNews(w, r)
	case http.MethodGet:
		fmt.Println("Get request!")
		h.GetNewsList(w, r)
	}
}

func (h *NewsHandler) CreateNews(w http.ResponseWriter, r *http.Request) {
	var post DisciplinePost
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := utils.CreateFolder(h.path + post.Name)
	if result {
		fmt.Println("Discipline created!")
	} else {
		fmt.Println("Failed create discipline!")
	}
}

func (h *NewsHandler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	var delete DisciplineDelete
	err := json.NewDecoder(r.Body).Decode(&delete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := os.Remove(h.path + delete.Name)
	if err != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("Discipline deleted!")
}

func (h *NewsHandler) PutNews(w http.ResponseWriter, r *http.Request) {
	var put DisciplinePut
	err := json.NewDecoder(r.Body).Decode(&put)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(put)

	err2 := os.Rename(h.path+put.OldName, h.path+put.NewName)
	if err != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("Discipline renamed!")
}

func (h *NewsHandler) GetNewsList(w http.ResponseWriter, r *http.Request) {
	news := []*News{}
	// bson.M{} - это типа условия для поиска
	err := h.News.Find(bson.M{}).All(&news)
	if err != nil {
		panic(err)
	}

	jsonNews, err := json.Marshal(news)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	ResponseWithJSON(w, jsonNews)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(json)
}
