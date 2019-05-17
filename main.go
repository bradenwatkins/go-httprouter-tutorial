package main

import (
	"log"
	"net/http"
)

func main() {

	store := &InMemoryBookStore{}

	router := NewRouter(store)
	log.Fatal(http.ListenAndServe(":8080", router))
}
