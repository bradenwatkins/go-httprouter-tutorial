package main

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Logger(fn func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		fn(w, r, param)
		log.Printf("Done in %v (%s %s)", time.Since(start), r.Method, r.URL.Path)
	}
}

func PanicLogger(fn func(w http.ResponseWriter, r *http.Request, param interface{})) func(w http.ResponseWriter, r *http.Request, param interface{}) {
	return func(w http.ResponseWriter, r *http.Request, param interface{}) {
		log.Panicf("SERVER ERROR:%s %s %v", r.Method, r.URL.Path, time.Now())
		fn(w, r, param)
	}
}
