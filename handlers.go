package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	writeOKResponse(w, fmt.Sprint("Welcome!"))
}

func BookCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	book := &Book{}
	if err := populateModelFromHandler(w, r, params, book); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	bookstore[book.ISDN] = book
	writeOKResponse(w, book)
}

func BookIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	books := make([]*Book, len(bookstore))
	idx := 0
	for _, book := range bookstore {
		books[idx] = book
		idx++
	}
	writeOKResponse(w, books)
}

func BookShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isdn := params.ByName("isdn")
	book, ok := bookstore[isdn]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if !ok {
		writeErrorResponse(w, http.StatusNotFound, "RecordNotFound")
	} else {
		writeOKResponse(w, book)
	}
}

func ServerErrorHandler(w http.ResponseWriter, r *http.Request, params interface{}) {
	writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
}

func writeOKResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&JsonResponse{Data: data})
}

func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})
}

func populateModelFromHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<32))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}
