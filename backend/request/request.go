package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DeleteRequest struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

func FindWithPagination(c *mgo.Collection, skip int, limit int, e interface{}) error {
	// bson.M{} - это типа условия для поиска
	err := c.Find(bson.M{}).Skip(skip).Limit(limit).All(e)
	return err
}

func ParseSkipAndLimit(str string) (skip, limit int) {
	regex, _ := regexp.Compile("\\d+")
	strings := regex.FindAllString(str, -1)

	skip, _ = strconv.Atoi(strings[0])
	end, _ := strconv.Atoi(strings[1])
	limit = end - skip + 1
	return
}

func GetWithPagination(c *mgo.Collection, rangeStr string, list interface{}) {
	skip, limit := ParseSkipAndLimit(rangeStr)
	err := FindWithPagination(c, skip, limit, list)
	// bson.M{} - это типа условия для поиска
	if err != nil {
		panic(err)
	}
}

func GetList(w http.ResponseWriter, c *mgo.Collection, rangeStr string, list interface{}) {
	fmt.Println("Get some list request!")

	if rangeStr != "" {
		GetWithPagination(c, rangeStr, list)
	} else {
		err := c.Find(bson.M{}).All(list)
		if err != nil {
			panic(err)
		}
	}

	count, err := c.Count()
	if err != nil {
		panic(err)
	}

	jsonNews, err := json.Marshal(list)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Get some list success!")

	contentRange := fmt.Sprint("news */", count)
	w.Header().Set("Content-Range", contentRange)
	ResponseWithJSON(w, jsonNews)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(json)
}

func GetOne(w http.ResponseWriter, c *mgo.Collection, id string, v interface{}) {
	err := c.FindId(bson.ObjectIdHex(id)).One(v)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(v)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Get some one success!")

	ResponseWithJSON(w, json)
}

func Decode(w http.ResponseWriter, body io.ReadCloser, v interface{}) {
	err := json.NewDecoder(body).Decode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func CreateOne(w http.ResponseWriter, c *mgo.Collection, v interface{}) {
	err2 := c.Insert(v)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(v)
	if err != nil {
		log.Err(err).Msg("json marshall err")
		return
	}

	fmt.Println("Create some success!")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func DeleteById(w http.ResponseWriter, r *http.Request, c *mgo.Collection) {
	id := chi.URLParam(r, "id")
	fmt.Println("Delete one news request! id: ", id)

	// TODO проверка на bson id

	err2 := c.RemoveId(bson.ObjectIdHex(id))
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Delete some success!")
}
type CRUD interface {
	GetList()
	GetOne()
	Update()
	Delete()
}