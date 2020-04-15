package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	database := session.DB("neurotech")

	publicationCollection := initPublicationsCollection(database)
	publicationHandler := &PublicationHandler{Publication: publicationCollection}
	http.Handle("/publications", publicationHandler)

	newsCollection := initNewsCollection(database)
	newsHandler := &NewsHandler{News: newsCollection}
	http.Handle("/news", newsHandler)

	disciplineHandler := &DisciplineHandler{path: DisciplinePath}
	http.Handle("/disciplines", disciplineHandler)

	studyHandler := &StudyMaterialHandler{path: DisciplinePath}
	http.Handle("/study-materials", studyHandler)

	http.ListenAndServe(":8080", nil)
}

func initPublicationsCollection(database *mgo.Database) *mgo.Collection {
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

func initNewsCollection(database *mgo.Database) *mgo.Collection {
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
