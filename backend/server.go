package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type News struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Title   string        `json:"title" bson:"title"`
	Content string        `json:"content" bson:"content"`
}

func main() {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// если коллекции не будет, то она создасться автоматически
	collection := session.DB("neurotech").C("news")

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

	newsHandler := &NewsHandler{
		path: NewsPath,
		News: collection,
	}
	http.Handle("/news", newsHandler)

	disciplineHandler := &DisciplineHandler{path: DisciplinePath}
	http.Handle("/disciplines", disciplineHandler)

	studyHandler := &StudyMaterialHandler{path: DisciplinePath}
	http.Handle("/study-materials", studyHandler)

	http.ListenAndServe(":8080", nil)
}
