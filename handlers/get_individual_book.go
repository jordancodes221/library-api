package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetIndividualBook allows the client to get an individual book in the library by its ISBN
func (h *BooksHandler) GetIndividualBook(c *gin.Context) {
	isbn := c.Param("isbn")
	book, err := h.BookDAOInterface.Read(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code if unsuccessful
		return
	}

	if book == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"details": "REQUEST SUCCESSFUL. BOOK NOT FOUND"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}