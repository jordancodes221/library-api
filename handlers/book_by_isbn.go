package handlers

import (
	"example/library_project/models"
)

// bookByISBN takes an ISBN string as inuput and returns the *Book from our library with that ISBN
func (h *BooksHandler) bookByISBN(isbn string) (*models.Book, error) {
	bookPtr, ok := h.Books[isbn] // in the future, this could be a call to a database

	// For scalability, we can add a database connection here. 
	// If there is an error connecting to the database, then we will return: nil, InternalServerError

	if ok {
		return bookPtr, nil
	} else {
		return nil, nil
	}
}