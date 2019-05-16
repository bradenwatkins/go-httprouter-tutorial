package main

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func AllRoutes() Routes {
	routes := Routes{
		Route{"Index", "GET", "/", Index},
		Route{"BookIndex", "GET", "/books", BookIndex},
		Route{"Bookshow", "GET", "/books/:isdn", BookShow},
		Route{"Bookshow", "POST", "/books", BookCreate},
	}
	return routes
}
