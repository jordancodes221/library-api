package handlers

import (	
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetAllBooks allows the client to get all of the books in the library
func (h *BooksHandler) GetAllBooks(c *gin.Context) {
	all_books, err := h.BookDAOInterface.ReadAll()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code if unsuccessful
		return
	}

	c.IndentedJSON(http.StatusOK, all_books)
}