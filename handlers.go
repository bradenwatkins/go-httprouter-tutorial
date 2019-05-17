package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/julienschmidt/httprouter"
)

func (router *Router) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	writeOKResponse(w, fmt.Sprint("Welcome!"))
}

func (router *Router) BookCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	book := &Book{}
	if err := populateModelFromHandler(w, r, params, book); err != nil {
		println(err.Error())
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	_, err := router.Store.Add(r.Context(), book)
	if err != nil {
		println(err.Error())
		writeErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeOKResponse(w, book)
}

func (router *Router) BookIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	books, err := router.Store.GetAll(r.Context())
	if err != nil {
		println(err.Error())
		writeErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeOKResponse(w, books)
}

func (router *Router) BookShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isdn := params.ByName("isdn")
	book, err := router.Store.Get(r.Context(), isdn)

	if xerrors.Is(err, BookNotFound) {
		writeErrorResponse(w, http.StatusNotFound, "book not found")
		return
	} else if err != nil {
		println(err.Error())
		writeErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeOKResponse(w, book)
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
