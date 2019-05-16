package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestBookIndex(t *testing.T) {
	testBook := &Book{
		ISDN:   "111",
		Title:  "test title",
		Author: "test author",
		Pages:  42,
	}
	bookstore[testBook.ISDN] = testBook

	req := newRequest(t, "GET", "/books", nil)
	rr := newRequestRecorder(req, "GET", "/books", BookIndex)
	expectedBody := "{\"meta\":null,\"data\":[{\"isdn\":\"111\",\"title\":\"test title\",\"author\":\"test author\",\"pages\":42}]}\n"
	checkResponseCode(t, http.StatusOK, rr.Code)
	checkResponseBody(t, expectedBody, rr.Body.String())
}

func TestBookCreate(t *testing.T) {
	testBook := &Book{
		ISDN:   "111",
		Title:  "test title",
		Author: "test author",
		Pages:  42,
	}

	req := newRequest(t, "PUT", "/books", testBook)
	rr := newRequestRecorder(req, "PUT", "/books", BookCreate)
	expected := "{\"meta\":null,\"data\":{\"isdn\":\"111\",\"title\":\"test title\",\"author\":\"test author\",\"pages\":42}}\n"
	checkResponseCode(t, http.StatusOK, rr.Code)
	checkResponseBody(t, expected, rr.Body.String())
}

func TestBookShow(t *testing.T) {
	testBook := &Book{
		ISDN:   "111",
		Title:  "test title",
		Author: "test author",
		Pages:  42,
	}
	bookstore[testBook.ISDN] = testBook

	// Test a GET call for a record that does exist
	req1 := newRequest(t, "GET", "/books/111", nil)
	rr1 := newRequestRecorder(req1, "GET", "/books/:isdn", BookShow)
	expected1 := "{\"meta\":null,\"data\":{\"isdn\":\"111\",\"title\":\"test title\",\"author\":\"test author\",\"pages\":42}}\n"
	checkResponseCode(t, http.StatusOK, rr1.Code)
	checkResponseBody(t, expected1, rr1.Body.String())

	// Test a GET call for a record that does NOT exist
	req2 := newRequest(t, "GET", "/books/222", nil)
	rr2 := newRequestRecorder(req2, "GET", "/books/:isdn", BookShow)
	expected2 := "{\"error\":{\"status\":404,\"title\":\"RecordNotFound\"}}\n"
	checkResponseCode(t, http.StatusNotFound, rr2.Code)
	checkResponseBody(t, expected2, rr2.Body.String())
}

func newRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	encodedBody, err := json.Marshal(body)
	req, err := http.NewRequest(method, path, strings.NewReader(string(encodedBody)))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func newRequestRecorder(req *http.Request, method string, strPath string, fnHandler func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.Handle(method, strPath, fnHandler)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Error(formatError("Response Codes Differ", expected, actual))
	}
}

func checkResponseBody(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Error(formatError("Response Bodies Differ", expected, actual))
	}
}

func formatError(message string, expected, actual interface{}) string {
	return fmt.Sprintf("%s\n%10s%s\n%10s%s\n", message, "expected:", fmt.Sprint(expected), "actual:", fmt.Sprint(actual))
}
