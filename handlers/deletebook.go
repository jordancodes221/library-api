package handlers

import ( 
	// "example/library_project/handlers"
	// "example/library_project/validators"
	// "example/library_project/models"
	


	"net/http"
	"github.com/gin-gonic/gin"
	// "errors"
	// "time"
	// "encoding/json"
	// "fmt"
	// "reflect"
	// "strconv"
)

// DELETE
func DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")

	book, err := bookByISBN(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code
		return
	}

	if book == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"details": "REQUEST SUCCESSFUL. BOOK NOT FOUND"})
		return
	}

	delete(mapOfBooks, isbn)
	c.Status(http.StatusNoContent) // 204 status code if successful
}