package main

import (
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/volatiletech/authboss/auth"
)


func debugln(args ...interface{}) {
	fmt.Println(args...)
}

func debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		h.ServeHTTP(w, r)
	})
}