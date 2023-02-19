package handlers

import ( // h.Books, bookByISBN
	"example/library_project/models"
	
	"net/http"
	"github.com/gin-gonic/gin"
)

// GET (all books)
func (h *BooksHandler) GetAllBooks(c *gin.Context) {
	// Make a slice containing all the values from mapOfBooks
	var vals []*models.Book
	
	for _, v := range h.Books { // should change mapOfBooks to h.Books
		vals = append(vals, v)
	}

	c.IndentedJSON(http.StatusOK, vals)
}

// GET (individual book)
func (h *BooksHandler) GetIndividualBook(c *gin.Context) {
	isbn := c.Param("isbn")
	book, err := h.bookByISBN(isbn)

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