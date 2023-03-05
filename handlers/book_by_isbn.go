package handlers

import ( // h.Books
	"example/library_project/models"
)

// Helper function
func (h *BooksHandler) bookByISBN(isbn string) (*models.Book, error) {
	bookPtr, ok := h.Books[isbn] // in the future, this could be a call to a database
	// if there is an error connecting to the database, then we will return: nil, InternalServerError

	if ok {
		return bookPtr, nil
	} else {
		return nil, nil
	}
}