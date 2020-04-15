package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	database := session.DB("neurotech")

	studenWorkCollection := InitStudentWorksCollection(database)
	studenWorkHandler := &StudentWorkHandler{Collection: studenWorkCollection}
	http.Handle("/student-work", studenWorkHandler)

	publicationCollection := InitPublicationsCollection(database)
	publicationHandler := &PublicationHandler{Collection: publicationCollection}
	http.Handle("/publications", publicationHandler)

	newsCollection := InitNewsCollection(database)
	newsHandler := &NewsHandler{Collection: newsCollection}
	http.Handle("/news", newsHandler)

	disciplineHandler := &DisciplineHandler{path: DisciplinePath}
	http.Handle("/disciplines", disciplineHandler)

	studyHandler := &StudyMaterialHandler{path: DisciplinePath}
	http.Handle("/study-materials", studyHandler)

	http.ListenAndServe(":8080", nil)
}
