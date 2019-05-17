package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"
)

// func TestBookIndex(t *testing.T) {
// 	testBook := &Book{
// 		ISDN:   "111",
// 		Title:  "test title",
// 		Author: "test author",
// 		Pages:  42,
// 	}
// 	bookstore[testBook.ISDN] = testBook

// 	req := newRequest(t, "GET", "/books", nil)
// 	rr := newRequestRecorder(req, "GET", "/books", BookIndex)
// 	expectedBody := "{\"meta\":null,\"data\":[{\"isdn\":\"111\",\"title\":\"test title\",\"author\":\"test author\",\"pages\":42}]}\n"
// 	checkResponseCode(t, http.StatusOK, rr.Code)
// 	checkResponseBody(t, expectedBody, rr.Body.String())
// }

// func TestBookCreate(t *testing.T) {
// 	testBook := &Book{
// 		ISDN:   "111",
// 		Title:  "test title",
// 		Author: "test author",
// 		Pages:  42,
// 	}

// 	req := newRequest(t, "PUT", "/books", testBook)
// 	rr := newRequestRecorder(req, "PUT", "/books", BookCreate)
// 	expected := "{\"meta\":null,\"data\":{\"isdn\":\"111\",\"title\":\"test title\",\"author\":\"test author\",\"pages\":42}}\n"
// 	checkResponseCode(t, http.StatusOK, rr.Code)
// 	checkResponseBody(t, expected, rr.Body.String())
// }

func TestBookShow(t *testing.T) {
	sampleBook1 := &Book{
		ISDN:   "111",
		Title:  "test title",
		Author: "test author",
		Pages:  42,
	}

	var tests = []struct {
		name   string
		book   *Book
		status int
	}{
		{"will 404 when missing a book", nil, http.StatusNotFound},
		{"will 200 with book when found", sampleBook1, http.StatusOK},
	}

	url := "/books/1"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &MockBookStore{
				get: func(ctx context.Context, id string) (*Book, error) {
					if test.book == nil {
						return nil, BookNotFound
					}
					return sampleBook1, nil
				},
			}

			router := NewRouter(store)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", url, nil)

			router.ServeHTTP(rr, r.WithContext(context.Background()))

			checkResponseCode(t, test.status, rr.Code)
			// checkResponseBody(t, "", rr.Body.String())
		})
	}
}

// func newRequest(t *testing.T, method, path string, body interface{}) *http.Request {
// 	encodedBody, err := json.Marshal(body)
// 	req, err := http.NewRequest(method, path, strings.NewReader(string(encodedBody)))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return req
// }

// func newRequestRecorder(req *http.Request, method string, strPath string, fnHandler func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) *httptest.ResponseRecorder {
// 	router := httprouter.New()
// 	router.Handle(method, strPath, fnHandler)

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)
// 	return rr
// }

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

type MockBookStore struct {
	add    func(context.Context, *Book) (string, error)
	get    func(context.Context, string) (*Book, error)
	getAll func(context.Context) ([]*Book, error)
}

func (s *MockBookStore) Add(ctx context.Context, book *Book) (string, error) {
	return s.add(ctx, book)
}
func (s *MockBookStore) Get(ctx context.Context, id string) (*Book, error) {
	return s.get(ctx, id)
}
func (s *MockBookStore) GetAll(ctx context.Context) ([]*Book, error) {
	return s.getAll(ctx)
}
