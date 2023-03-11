package handlers

import (
	"example/library_project/models"
	
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetAllBooks allows the client to get all of the books in the library
func (h *BooksHandler) GetAllBooks(c *gin.Context) {
	// Make a slice containing all the values from mapOfBooks
	var vals []*models.Book
	
	for _, v := range h.Books { // should change mapOfBooks to h.Books
		vals = append(vals, v)
	}

	c.IndentedJSON(http.StatusOK, vals)
}