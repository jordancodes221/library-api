package handlers

import (
	"example/library_project/models"
	"example/library_project/utils"
	"example/library_project/dao"
)

// BooksHandlers is the struct on which all handler functions are defined as pointer-receiver functions
type BooksHandler struct {
	// Books is the library of all the books
    Books map[string]*models.Book
	DateTimeInterface utils.DateTimeProvider
}

func (h *BooksHandler) Create(newBook *models.Book) {
	h.Books[*newBook.ISBN] = newBook
}

func (h *BooksHandler) Delete(book *models.Book) {
	delete(h.Books, *book.ISBN)
}

func (h *BooksHandler) Update(book *models.Book) {
	h.Books[*book.ISBN] = book
}

func (h *BooksHandler) Read(isbn string) (*models.Book, error) {
	retrievedBook, ok := h.Books[isbn] // in the future, this could be a call to a database

	// For scalability, we can add a database connection here. 
	// If there is an error connecting to the database, then we will return: nil, InternalServerError

	if ok {
		return retrievedBook, nil
	} else {
		return nil, nil
	}
}

func (h *BooksHandler) ReadAll() ([]*models.Book, error) {	
	all_books := make([]*models.Book, 0)

	// For scalability, we can add a database connection here. 
	// If there is an error connecting to the database, then we will return: nil, InternalServerError

	for _, currentBook := range h.Books {
		all_books = append(all_books, currentBook)
	}

	return all_books, nil
}

