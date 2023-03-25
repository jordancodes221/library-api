package handlers

import (
	"example/library_project/models"
	"example/library_project/utils"

	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"encoding/json"
	// "errors"
)

// CreateBook allows the client to add a new book to the library
func (h *BooksHandler) CreateBook(c *gin.Context) {
	// Decode JSON to book struct
	newBook := new(models.Book) // the "new" keyword allocates memory for models.Book, and returns a pointer to it
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// If fields are not nil, ensure they are within range
	if err := newBook.Validate(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Logic validation
	if err := newBook.ValidateLogicForCreateBook(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Make sure ISBN is not already in-use
	if _, ok := h.Books[*newBook.ISBN]; ok {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": "Book already exists."})
		return
	}

	// Update TimeCreated to now
	newBook.TimeCreated = utils.ToPtr(time.Now())

	// Add the new book to our library
	h.Books[*newBook.ISBN] = newBook

	c.IndentedJSON(http.StatusCreated, newBook) // 201 status code if successful
}