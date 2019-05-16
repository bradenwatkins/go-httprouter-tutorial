package main

import (
	"github.com/julienschmidt/httprouter"
)

func NewRouter(routes Routes) *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, Logger(route.HandlerFunc))
	}
	router.PanicHandler = PanicLogger(ServerErrorHandler)
	return router
}
