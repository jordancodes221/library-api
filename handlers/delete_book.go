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

	h.Delete(book)
	c.Status(http.StatusNoContent)
}