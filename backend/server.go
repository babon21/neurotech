package main

import (
	"net/http"
)

func main() {
	disciplineHandler := &DisciplineHandler{path: DisciplinePath}
	http.Handle("/disciplines", disciplineHandler)

	studyHandler := &StudyMaterialHandler{path: DisciplinePath}
	http.Handle("/study-materials", studyHandler)

	http.ListenAndServe(":8080", nil)
}
