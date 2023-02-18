package handlers

import ( 
	// "example/library_project/handlers"
	// "example/library_project/validators"
	"example/library_project/models"
	


	// "net/http"
	// "github.com/gin-gonic/gin"
	// "errors"
	"time"
	// "encoding/json"
	// "fmt"
	// "reflect"
	// "strconv"
)

// Generic function converts literals to pointers
func ToPtr[T string|time.Time](v T) *T {
    return &v
}

// Helper function
func bookByISBN(isbn string) (*models.Book, error) {
	bookPtr, ok := mapOfBooks[isbn] // in the future, this could be a call to a database
	// if there is an error connecting to the database, then we will return: nil, InternalServerError

	if ok {
		return bookPtr, nil
	} else {
		return nil, nil
	}
}