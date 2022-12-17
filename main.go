package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	// "fmt"
	// "encoding/json"
	// "strconv"
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
var bookInstance0 Book = Book{"0000", "available", 	"", 	"", 		time.Now().String(), time.Time{}.String()} // --> Available
var bookInstance1 Book = Book{"0001", "available", 	"", 	"", 		time.Now().String(), time.Time{}.String()} // --> Checked-out
var bookInstance2 Book = Book{"0002", "available", 	"", 	"", 		time.Now().String(), time.Time{}.String()} // --> On-hold

var bookInstance3 Book = Book{"0003", "checked-out", 	"", 	"01", 		time.Now().String(), time.Time{}.String()} // --> Available
var bookInstance4 Book = Book{"0004", "checked-out", 	"", 	"01", 		time.Now().String(), time.Time{}.String()} // --> Available (no match)
var bookInstance5 Book = Book{"0005", "checked-out", 	"", 	"01", 		time.Now().String(), time.Time{}.String()} // --> Checked-out
var bookInstance6 Book = Book{"0006", "checked-out", 	"", 	"01", 		time.Now().String(), time.Time{}.String()} // --> Checked-out (no match)
var bookInstance7 Book = Book{"0007", "checked-out", 	"", 	"01", 		time.Now().String(), time.Time{}.String()} // --> On-hold 
var bookInstance8 Book = Book{"0008", "checked-out", 	"", 	"01", 		time.Now().String(), time.Time{}.String()} // --> On-hold (no match)

var bookInstance9 Book =  Book{"0009", "on-hold", 	"01", 	"", 		time.Now().String(), time.Time{}.String()} // --> Available
var bookInstance10 Book = Book{"0010", "on-hold", 	"01", 	"", 		time.Now().String(), time.Time{}.String()} // --> Available (no match)
var bookInstance11 Book = Book{"0011", "on-hold", 	"01", 	"", 		time.Now().String(), time.Time{}.String()} // --> Checked-out
var bookInstance12 Book = Book{"0012", "on-hold", 	"01", 	"", 		time.Now().String(), time.Time{}.String()} // --> Checked-out (no match)
var bookInstance13 Book = Book{"0013", "on-hold", 	"01", 	"", 		time.Now().String(), time.Time{}.String()} // --> On-hold 
var bookInstance14 Book = Book{"0014", "on-hold", 	"01", 	"", 		time.Now().String(), time.Time{}.String()} // --> On-hold (no match)

var bookInstance15 Book = Book{"0015", "available", 	"", 	"", 		time.Now().String(), time.Time{}.String()} // --> This is the book to be deleted in testing

var mapOfBooks = map[string]*Book{
	"0000" : &bookInstance0,
	"0001" : &bookInstance1,
	"0002" : &bookInstance2,

	"0003" : &bookInstance3,
	"0004" : &bookInstance4,
	"0005" : &bookInstance5,
	"0006" : &bookInstance6,
	"0007" : &bookInstance7,
	"0008" : &bookInstance8,

	"0009" : &bookInstance9,
	"0010" : &bookInstance10,
	"0011" : &bookInstance11,
	"0012" : &bookInstance12,
	"0013" : &bookInstance13,
	"0014" : &bookInstance14,

	"0015" : &bookInstance15,
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// Input checking
func inputCheck(book *Book) error {
	// Consider adding something to check ISBN number format (see: https://pkg.go.dev/github.com/moraes/isbn)

	// check if state is valid
	if (book.State != "available") && (book.State != "on-hold") && (book.State != "checked-out") {
		return errors.New("Invalid state in created book.")
	}

	// if book is available, it cannot have an on-hold or checked-out customer ID
	if (book.State == "available") {
		if (book.OnHoldCustomerID != "") && (book.CheckedOutCustomerID == "") {
			return errors.New("Cannot have an on-hold customer ID on an available book.")
		}

		if (book.OnHoldCustomerID == "") && (book.CheckedOutCustomerID != "") {
			return errors.New("Cannot have checked-out customer ID on an available book.")
		}

		if (book.OnHoldCustomerID != "") && (book.CheckedOutCustomerID != "") {
			return errors.New("Cannot have an on-hold or checked-out customer ID on an available book.")
		}
	}

	// If new book is on-hold, ensure there is an on-hold customer ID and no checked-out ID
	if (book.State == "on-hold") {
		if (book.OnHoldCustomerID == "") {
			return errors.New("Missing on-hold customer ID.")
		}

		if (book.CheckedOutCustomerID != "") {
			return errors.New("Cannot have a checked-out customer ID on an on-hold book.")
		}
	}

	// If new book is checked-out, ensure there is a checked-out customer ID and no on-hold ID
	if (book.State == "checked-out") {
		if (book.CheckedOutCustomerID == "") {
			return errors.New("Missing checked-out customer ID.")
		}

		if (book.OnHoldCustomerID != "") {
			return errors.New("Cannot have an on-hold customer ID on a checked-out book.")
		}
	}

	return nil
}

// POST
func CreateBook(c *gin.Context) {
	var newBook Book

	// Unmarshal
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Make sure ISBN is not already in-use
	if _, ok := mapOfBooks[newBook.ISBN]; ok {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": "Book already exists."})
		return
	}

	// Input checking
	if err := inputCheck(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	newBook.TimeCreated = time.Now().String()
	newBook.TimeUpdated = time.Now().String()

	mapOfBooks[newBook.ISBN] = &newBook

	c.IndentedJSON(http.StatusCreated, newBook) // 201 status code if successful
}

// DELETE
func DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")

	_, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": err.Error()})
		return
	}

	delete(mapOfBooks, isbn)
	c.Status(http.StatusNoContent) // 204 status code if successful
}

var NoMatchError error = errors.New("ID's do not match.")

// Checkout
	// available --> checked-out
	// on-hold --> checked-out
	// checked-out --> checked-out
func Checkout(currentBook *Book, incomingBook *Book) (*Book, error) {
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
			return nil, NoMatchError
		}
	}

	if (currentBook.State == "checked-out") {
		if (currentBook.CheckedOutCustomerID == incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it checked out is the same one trying to check it out redundantly
			// pass
		} else {
			return nil, NoMatchError
		}
	}

	return currentBook, nil
}

// Conflict
	// checked-out --> on-hold
func Conflict(currentBook *Book, incomingBook *Book) (*Book, error) {
	return nil, errors.New("Invalid state transfer requested.")
}

// PlaceHold
	// available --> on-hold
	// on-hold --> on-hold
func PlaceHold(currentBook *Book, incomingBook *Book) (*Book, error) {
	if (currentBook.State == "available") {
		currentBook.State = "on-hold"
		currentBook.OnHoldCustomerID = incomingBook.OnHoldCustomerID
		currentBook.TimeUpdated = time.Now().String()
	}

	if (currentBook.State == "on-hold") {
		if (currentBook.OnHoldCustomerID == incomingBook.OnHoldCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			// pass
		} else {
			return nil, NoMatchError
		}
	}

	return currentBook, nil
}

// ReleaseHold
	// on-hold --> available (when ID's match)
func ReleaseHold(currentBook *Book, incomingBook *Book) (*Book, error) {
	if (currentBook.State == "on-hold") {
		if (currentBook.OnHoldCustomerID == incomingBook.OnHoldCustomerID) {
			currentBook.State = "available"
			currentBook.OnHoldCustomerID = ""
			currentBook.TimeUpdated = time.Now().String()
		} else {
			return nil, NoMatchError
		}
	}

	return currentBook, nil
}

// Return
	// checked-out --> available (when ID's match)
func Return(currentBook *Book, incomingBook *Book) (*Book, error) {
	if (currentBook.State == "checked-out") {
		if (currentBook.CheckedOutCustomerID == incomingBook.CheckedOutCustomerID) {
			currentBook.State = "available"
			currentBook.CheckedOutCustomerID = "" // need this, or leave it as who most recently had it on hold?
			currentBook.TimeUpdated = time.Now().String()
		} else {
			return nil, NoMatchError
		}
	}

	return currentBook, nil
}

// NoOperation
	// available --> available
	// on-hold --> on-hold (when ID's match)
func NoOperation(currentBook *Book, incomingBook *Book) (*Book, error) {
	return currentBook, nil
}

// First key is current state, 2nd key is incoming state
var actionTable = map[string]map[string]func(currentBook *Book, incomingBook *Book) (*Book, error) {
	"available": {
		"available": NoOperation,
		"checked-out": Checkout,
		"on-hold": PlaceHold,
	}, "checked-out": {
			"available": Return,
			"checked-out": Checkout,
			"on-hold": Conflict,
	}, "on-hold": {
			"available": ReleaseHold,
			"checked-out": Checkout,
			"on-hold": PlaceHold,
	},
}

// PATCH
func UpdateBook(c *gin.Context) {
	isbn := c.Param("isbn")

	// Ensure book to be updated exists
	book, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": err.Error()})
		return
	}

	// Unmarshal JSON
	var incomingRequest *Book
	if err := c.BindJSON(&incomingRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Create variables for the 2 states - current and incoming
	currentState := book.State
	incomingState := incomingRequest.State

	// Input checking
	if (incomingState != "available") && (incomingState != "on-hold") && (incomingState != "checked-out") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Invalid state in incoming request."})
		return
	}

	// Call the appropriate function from  the action table, and catch possible errors
	book, err = actionTable[currentState][incomingState](book, incomingRequest)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

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
		// curl localhost:8080/books/0005 --request "DELETE"
	// PATCH
		// curl -X PATCH localhost:8080/books/0001 -H 'Content-Type: application/json' -H 'Accept: application/json' -d @incomingRequest.json