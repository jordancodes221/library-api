package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	// "encoding/json"
	// "strconv"
	// "fmt"
)

type Book struct{
	ISBN 					string 	`json:"isbn"`
	State 					string 	`json:"state"`

	OnHoldCustomerID 		string 	`json:"onholdcustomerid"` 	// re-named from Onhold_customerID
	CheckedOutCustomerID 	string 	`json:"checkedoutcustomerid"` // re-named from CheckedOutCustomerID

	TimeCreated 			string 	`json:"timecreated"` // re-named from Time_created
	TimeUpdated  			string 	`json:"timeupdated"` // re-named from Time_updated
}

// Test data
var bookInstance0 Book = Book{"0000", "available", 	"", 	"", 		time.Now().String(), time.Time{}.String()} // not on-hold, not checked-out
var bookInstance1 Book = Book{"0001", "checked-out", 	"", 	"01", 	time.Now().String(), time.Time{}.String()} // checked-out, not on-hold
var bookInstance2 Book = Book{"0002", "checked-out", 	"02", 	"01", 	time.Now().String(), time.Time{}.String()} // checked-out by one customer, on-hold by another

var mapOfBooks = map[string]*Book{
	"0000" : &bookInstance0,
	"0001" : &bookInstance1,
	"0002" : &bookInstance2,
}

// GET (all books)
func GetAllBooks(c *gin.Context) {
	// Make a slice containing all the values from mapOfBooks
	var vals []*Book
	
	for _, v := range mapOfBooks {
		vals = append(vals, v)
	}

	c.IndentedJSON(http.StatusOK, vals)
}

// Helper function for GET (individual book)
func bookByISBN(isbn string) (*Book, error) {
	bookPtr, ok := mapOfBooks[isbn]

	if ok {
		return bookPtr, nil
	} else {
		return nil, errors.New("Book not found.")
	}
}

// GET (individual book)
func GetIndividualBook(c *gin.Context) {
	isbn := c.Param("isbn")
	book, err := bookByISBN(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// POST
func CreateBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		return // BindJSON handles the error response
	}

	newBook.TimeCreated = time.Now().String()
	newBook.TimeUpdated = time.Now().String()

	mapOfBooks[newBook.ISBN] = &newBook

	c.IndentedJSON(http.StatusCreated, newBook) // 201 status code if successful
}

// DELETE
func DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")
	delete(mapOfBooks, isbn)
	c.Status(http.StatusNoContent) // 204 status code if successful
}

// Checkout
	// available --> checked-out
	// on-hold --> checked-out
	// checked-out --> checked-out
func Checkout(c *gin.Context, currentBook *Book, incomingBook *Book) {
	if (currentBook.State == "available") {
		currentBook.State = "checked-out" // or should we use incomingBook.State? 
		currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
		currentBook.TimeUpdated = time.Now().String()
	}

	if (currentBook.State == "on-hold") {
		if (currentBook.OnHoldCustomerID == incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			currentBook.State = "checked-out"
			currentBook.OnHoldCustomerID = ""
			currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
			currentBook.TimeUpdated = time.Now().String()
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "IDs do not match."})
		}
	}

	if (currentBook.State == "checked-out") {
		if (currentBook.CheckedOutCustomerID == incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it checked out is the same one trying to check it out redundantly
			// pass
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "IDs do not match."})
		}
	}
}

// Conflict
	// checked-out --> on-hold
func Conflict(c *gin.Context, currentBook *Book, incomingBook *Book) {
	c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Invalid state requested"})
}

// PlaceHold
	// available --> on-hold
	// on-hold --> on-hold
func PlaceHold(c *gin.Context, currentBook *Book, incomingBook *Book) {
	if (currentBook.State == "available") {
		currentBook.State = "on-hold"
		currentBook.OnHoldCustomerID = incomingBook.OnHoldCustomerID
		currentBook.TimeUpdated = time.Now().String()
	}

	if (currentBook.State == "on-hold") {
		if (currentBook.OnHoldCustomerID == incomingBook.OnHoldCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			// pass
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "IDs do not match."})
		}
	}
}

// ReleaseHold
	// on-hold --> available (when ID's match)
func ReleaseHold(c *gin.Context, currentBook *Book, incomingBook *Book) {
	if (currentBook.State == "on-hold") {
		if (currentBook.OnHoldCustomerID == incomingBook.OnHoldCustomerID) {
			currentBook.State = "available"
			currentBook.OnHoldCustomerID = "" // need this, or leave it as who most recently had it on hold?
			currentBook.TimeUpdated = time.Now().String()
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "IDs do not match."})
		}
	}
}

// Return
	// checked-out --> available (when ID's match)
func Return(c *gin.Context, currentBook *Book, incomingBook *Book) {
	if (currentBook.State == "checked-out") {
		if (currentBook.CheckedOutCustomerID == incomingBook.CheckedOutCustomerID) {
			currentBook.State = "available"
			currentBook.CheckedOutCustomerID = "" // need this, or leave it as who most recently had it on hold?
			currentBook.TimeUpdated = time.Now().String()
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "IDs do not match."})
		}
	}
}

// NoOperation
	// available --> available
	// on-hold --> on-hold (when ID's match)
func NoOperation(c *gin.Context, currentBook *Book, incomingBook *Book) {
	// pass
}

var actionTable = map[string]map[string]func(c *gin.Context, currentBook *Book, incomingBook *Book) {
	"available": {
		"available": NoOperation,
		"checked-out": Return,
		"on-hold": ReleaseHold,
	}, "checked-out": {
			"available": Checkout,
			"checked-out": NoOperation,
			"on-hold": Checkout,
	}, "on-hold": {
			"available": PlaceHold,
			"checked-out": Conflict,
			"on-hold": NoOperation,
	},
}

// PATCH
func UpdateBook(c *gin.Context) {
	isbn := c.Param("isbn")

	// This section could be removed if we change the resource PATCH acts on to be just /books, and we get the right book by the isbn in the incoming Book struct request
	book, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Book not found."})
		return
	}

	// Unmarshal JSON
	var incomingRequest *Book
	if err := c.BindJSON(&incomingRequest); err != nil {
		return
	}

	// Create variables for the 2 states - current and incoming
	currentState := book.State
	incomingState := incomingRequest.State

	// Input checking
	if (incomingState != "available") && (incomingState != "on-hold") && (incomingState != "checked-out") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Invalid state requested."})
		return
	}

	// Call the appropriate function from the action table
	actionTable[currentState][incomingState](c, book, incomingRequest)

	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()
	router.GET("/books", GetAllBooks)
	router.GET("/books/:isbn", GetIndividualBook)
	router.POST("/books", CreateBook)
	router.DELETE("/books/:isbn", DeleteBook)
	router.PATCH("/books/:isbn", UpdateBook)

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
		// curl localhost:8080/books/0000 --request "DELETE"
	// PATCH
		// curl -X PATCH localhost:8080/books/0000 -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"RequestedState": "checked-out", "CustomerID": "01"}'
		// curl -X PATCH localhost:8080/books/0000 -H 'Content-Type: application/json' -H 'Accept: application/json' -d @incomingRequest.json
			// in the 2nd command, we can change the endpoint because the individual book can be gotten from the isbn contained in the json file