package main

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/xerrors"
)

type Book struct {
	ISDN   string `json:"isdn"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

var BookNotFound = errors.New("book not found")

// InMemoryBookStore defines a book store that uses an in-memory map.
type InMemoryBookStore struct {
	mu    sync.RWMutex
	books map[string]*Book

	initOnce sync.Once
}

// Add will add a book to the in-memory map.
func (s *InMemoryBookStore) Add(ctx context.Context, book *Book) (string, error) {
	s.initOnce.Do(s.init)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.books[book.ISDN] = book

	return book.ISDN, nil
}

// Get will look up a book from the in-memory map or return a BookNotFound error.
func (s *InMemoryBookStore) Get(ctx context.Context, id string) (*Book, error) {
	s.initOnce.Do(s.init)

	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	if !ok {
		return nil, xerrors.Errorf(
			"no book with id: %s: %w", id, BookNotFound,
		)
	}

	return book, nil
}

// GetAll will return a list of all books in the in-memory map.
func (s *InMemoryBookStore) GetAll(ctx context.Context) ([]*Book, error) {
	s.initOnce.Do(s.init)

	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]*Book, len(s.books))
	idx := 0
	for _, book := range s.books {
		books[idx] = book
		idx++
	}

	return books, nil
}

// init does one-time setup
func (s *InMemoryBookStore) init() {

	if s.books == nil {
		s.books = make(map[string]*Book)
	}

}
