package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

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

	jsonNews, err := json.Marshal(news)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Create news success!")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonNews)
}

func (h *NewsHandler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete news request!")

	newsID := chi.URLParam(r, "newsID")
	fmt.Println("Delete one news request! id: ", newsID)

	// TODO проверка на bson id

	err2 := h.Collection.RemoveId(bson.ObjectIdHex(newsID))
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Delete news success!")
}

func (h *NewsHandler) UpdateNews(w http.ResponseWriter, r *http.Request) {
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

	rangeParam := r.URL.Query().Get("range")

	// need check to rangeParam
	news := []*News{}

	if rangeParam != "" {
		regex, _ := regexp.Compile("\\d+")
		strings := regex.FindAllString(rangeParam, -1)

		skip, _ := strconv.Atoi(strings[0])
		end, _ := strconv.Atoi(strings[1])
		limit := end - skip + 1

		// bson.M{} - это типа условия для поиска
		err := h.Collection.Find(bson.M{}).Skip(skip).Limit(limit).All(&news)
		if err != nil {
			panic(err)
		}
	} else {
		err := h.Collection.Find(bson.M{}).All(&news)
		if err != nil {
			panic(err)
		}
	}

	count, err := h.Collection.Count()
	if err != nil {
		panic(err)
	}

	jsonNews, err := json.Marshal(news)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Get news list success!")

	contentRange := fmt.Sprint("news */", count)
	w.Header().Set("Content-Range", contentRange)
	ResponseWithJSON(w, jsonNews)
}

func (h *NewsHandler) GetOneNews(w http.ResponseWriter, r *http.Request) {
	newsID := chi.URLParam(r, "newsID")
	fmt.Println("Get one news request! id: ", newsID)

	news := &News{}

	err := h.Collection.FindId(bson.ObjectIdHex(newsID)).One(&news)
	if err != nil {
		panic(err)
	}

	jsonNews, err := json.Marshal(news)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Get one news success!")

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
