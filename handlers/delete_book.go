package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// DeleteBook allows the client to delete a book from the library
func (h *BooksHandler) DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")

	book, err := h.bookByISBN(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code
		return
	}

	if book == nil {
		c.Status(http.StatusNoContent)
		return
	}

	delete(h.Books, isbn)
	c.Status(http.StatusNoContent)
}