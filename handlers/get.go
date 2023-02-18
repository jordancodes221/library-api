package handlers

import ( 
	"example/library_project/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

// GET (all books)
func GetAllBooks(c *gin.Context) {
	// Make a slice containing all the values from mapOfBooks
	var vals []*models.Book
	
	for _, v := range mapOfBooks {
		vals = append(vals, v)
	}

	c.IndentedJSON(http.StatusOK, vals)
}

// GET (individual book)
func GetIndividualBook(c *gin.Context) {
	isbn := c.Param("isbn")
	book, err := bookByISBN(isbn)

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