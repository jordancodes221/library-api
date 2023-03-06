package handlers

import ( 	 // h.Books, bookByISBN
	"net/http"
	"github.com/gin-gonic/gin"
)

// DELETE
func (h *BooksHandler) DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")

	book, err := h.bookByISBN(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code
		return
	}

	if book == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"details": "REQUEST SUCCESSFUL. BOOK NOT FOUND"})
		return
	}

	delete(h.Books, isbn)
	c.Status(http.StatusNoContent) // 204 status code if successful
}