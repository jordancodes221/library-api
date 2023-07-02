package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// DeleteBook allows the client to delete a book from the library
func (h *BooksHandler) DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")

	book, err := h.BookDAOInterface.Read(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code
		return
	}

	if book == nil {
		c.Status(http.StatusNoContent)
		return
	}

	if err = h.BookDAOInterface.Delete(book); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}