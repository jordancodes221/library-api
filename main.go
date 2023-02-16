package main

import ( 
	"example/library_project/handlers"

	// "example/library_project/validators"
	// "example/library_project/models"
	
	// "net/http"
	"github.com/gin-gonic/gin"
	// "errors"
	// "time"
	// "encoding/json"
	// "fmt"

	// "reflect"
	// "strconv"
)


func main() {
	router := gin.Default()
	router.GET("/books", handlers.GetAllBooks)
	router.GET("/books/:isbn", handlers.GetIndividualBook)
	router.POST("/books", handlers.CreateBook)
	router.DELETE("/books/:isbn", handlers.DeleteBook)
	router.PATCH("/books/:isbn", handlers.UpdateBook)

	router.Run("localhost:8080")
}

// To test, run "go run ." in one terminal window and a curl command in the another terminal window.
// Examples of curl commands are:
	// GET (all books)
		// curl localhost:8080/books
	// GET (individual book)
		// curl localhost:8080/books/0000
	// POST
		// curl localhost:8080/books --include --header "Content-Type: application/json" -d @newBook.json --request "POST"
	// DELETE
		// curl localhost:8080/books/0005 --request "DELETE"
	// PATCH
		// curl -X PATCH localhost:8080/books/00 -H 'Content-Type: application/json' -H 'Accept: application/json' -d @incomingRequest.json