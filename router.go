package main

import (
	"context"

	"github.com/julienschmidt/httprouter"
)

type BookStore interface {
	Add(context.Context, *Book) (string, error)
	Get(context.Context, string) (*Book, error)
	GetAll(context.Context) ([]*Book, error)
}

type Router struct {
	*httprouter.Router
	Store BookStore
}

func NewRouter(store BookStore) *Router {
	router := &Router{
		httprouter.New(),
		store,
	}

	router.GET("/", router.Index)
	router.GET("/books", router.BookIndex)
	router.GET("/books/:isdn", router.BookShow)
	router.POST("/books", router.BookCreate)

	router.PanicHandler = PanicLogger(ServerErrorHandler)
	return router
}
