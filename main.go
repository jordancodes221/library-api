package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	// "fmt"
	"encoding/json"
	// "reflect"
	// "strconv"
)

type Book struct{
	ISBN 			*string 	`json:"isbn"`
	State 			*string 	`json:"state"`

	OnHoldCustomerID 	*string 	`json:"onholdcustomerid"`
	CheckedOutCustomerID 	*string 	`json:"checkedoutcustomerid"`

	TimeCreated 		*time.Time 	`json:"timecreated"`
	TimeUpdated  		*time.Time	`json:"timeupdated"`
}

//// Instantiating Test Data
// Generic function converts literals to pointers
func ToPtr[T string|time.Time](v T) *T {
    return &v
}

// First test of instantiating test data with new schema and ToPtr function
var bookInstance00 Book = Book{ToPtr("00"), ToPtr("on-hold"), ToPtr("01"), nil, ToPtr(time.Now()), ToPtr(time.Now())}

// Actual test data to be used in testing
var bookInstance0 Book = Book{ToPtr("0000"), ToPtr("available"), 	nil, 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Available
var bookInstance1 Book = Book{ToPtr("0001"), ToPtr("available"), 	nil, 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Checked-out
var bookInstance2 Book = Book{ToPtr("0002"), ToPtr("available"), 	nil, 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> On-hold

var bookInstance3 Book = Book{ToPtr("0003"), ToPtr("checked-out"), 	nil, 	ToPtr("01"), 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Available
var bookInstance4 Book = Book{ToPtr("0004"), ToPtr("checked-out"), 	nil, 	ToPtr("01"), 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Available (no match)
var bookInstance5 Book = Book{ToPtr("0005"), ToPtr("checked-out"), 	nil, 	ToPtr("01"), 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Checked-out
var bookInstance6 Book = Book{ToPtr("0006"), ToPtr("checked-out"), 	nil, 	ToPtr("01"), 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Checked-out (no match)
var bookInstance7 Book = Book{ToPtr("0007"), ToPtr("checked-out"), 	nil, 	ToPtr("01"), 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> On-hold 
var bookInstance8 Book = Book{ToPtr("0008"), ToPtr("checked-out"), 	nil, 	ToPtr("01"), 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> On-hold (no match)

var bookInstance9 Book =  Book{ToPtr("0009"), ToPtr("on-hold"), 	ToPtr("01"), 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Available
var bookInstance10 Book = Book{ToPtr("0010"), ToPtr("on-hold"), 	ToPtr("01"), 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Available (no match)
var bookInstance11 Book = Book{ToPtr("0011"), ToPtr("on-hold"), 	ToPtr("01"), 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Checked-out
var bookInstance12 Book = Book{ToPtr("0012"), ToPtr("on-hold"), 	ToPtr("01"), 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> Checked-out (no match)
var bookInstance13 Book = Book{ToPtr("0013"), ToPtr("on-hold"), 	ToPtr("01"), 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> On-hold 
var bookInstance14 Book = Book{ToPtr("0014"), ToPtr("on-hold"), 	ToPtr("01"), 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> On-hold (no match)

var bookInstance15 Book = Book{ToPtr("0015"), ToPtr("available"), 	nil, 	nil, 	ToPtr(time.Now()), ToPtr(time.Time{})} // --> This is the book to be deleted in testing

// Map of test data to be used in testing
var mapOfBooks = map[string]*Book{
	"00" : &bookInstance00,

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
	bookPtr, ok := mapOfBooks[isbn] // in the future, this could be a call to a database
	// if there is an error connecting to the database, then we will return: nil, InternalServerError

	if ok {
		return bookPtr, nil
	} else {
		return nil, nil
	}
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"REQUEST SUCCESSFUL": "BOOK NOT FOUND"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// POST
func CreateBook(c *gin.Context) {
	var newBook *Book = &Book{nil, nil, nil, nil, nil, nil}

	// Unmarshal
	incomingBookAsMap := map[string]interface{}{}
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Update newBook with incoming ISBN and incoming State
	// ASSUME FOR NOW THAT ISBN AND STATE ARE ALWAYS PROVIDED
	incomingISBN := incomingBookAsMap["isbn"].(string)
	newBook.ISBN = ToPtr(incomingISBN)

	incomingState := incomingBookAsMap["state"].(string)
	newBook.State = ToPtr(incomingState)

	// Make sure ISBN is not already in-use
	if _, ok := mapOfBooks[incomingISBN]; ok {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": "Book already exists."})
		return
	}

	// Update newBook with incoming OnHoldCustomerID and incoming CheckedOutCustomerID, only if either if provided
	if incomingOnHoldCustomerID, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]; hasOnHoldCustomerID {
		incomingOnHoldCustomerID := incomingOnHoldCustomerID.(string) // Type assertion
		newBook.OnHoldCustomerID = ToPtr(incomingOnHoldCustomerID)
	}

	if incomingCheckedOutCustomerID, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]; hasCheckedOutCustomerID {
		incomingCheckedOutCustomerID := incomingCheckedOutCustomerID.(string) // Type assertion
		newBook.CheckedOutCustomerID = ToPtr(incomingCheckedOutCustomerID)
	}

	// Update newBook times
	newBook.TimeCreated = ToPtr(time.Now())
	newBook.TimeUpdated = ToPtr(time.Time{})

	// Add newBook to mapOfBooks
	mapOfBooks[*newBook.ISBN] = newBook

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
	if (*currentBook.State == "available") {
		*currentBook.State = "checked-out" // or should we use incomingBook.State? 
		currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
		currentBook.TimeUpdated = ToPtr(time.Now())
	} else if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			*currentBook.State = "checked-out"
			currentBook.OnHoldCustomerID = nil
			currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
			currentBook.TimeUpdated = ToPtr(time.Now())
		} else {
			return nil, NoMatchError
		}
	} else if (*currentBook.State == "checked-out") {
		if (*currentBook.CheckedOutCustomerID == *incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it checked out is the same one trying to check it out redundantly
			// pass
		} else {
			return nil, NoMatchError
		}
	} else {
		// pass
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
	if (*currentBook.State == "available") {
		*currentBook.State = "on-hold"
		currentBook.OnHoldCustomerID = incomingBook.OnHoldCustomerID
		currentBook.TimeUpdated = ToPtr(time.Now())
	} else if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.OnHoldCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			// pass
		} else {
			return nil, NoMatchError
		}
	} else {
		// pass 
	}

	return currentBook, nil
}

// ReleaseHold
	// on-hold --> available
func ReleaseHold(currentBook *Book, incomingBook *Book) (*Book, error) {
	if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.OnHoldCustomerID) {
			*currentBook.State = "available"
			currentBook.OnHoldCustomerID = nil
			currentBook.TimeUpdated = ToPtr(time.Now())
		} else {
			return nil, NoMatchError
		}
	}

	return currentBook, nil
}

// Return
	// checked-out --> available
func Return(currentBook *Book, incomingBook *Book) (*Book, error) {
	if (*currentBook.State == "checked-out") {
		if (*currentBook.CheckedOutCustomerID == *incomingBook.CheckedOutCustomerID) {
			*currentBook.State = "available"
			currentBook.CheckedOutCustomerID = nil
			currentBook.TimeUpdated = ToPtr(time.Now())
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

	incomingBookAsMap := map[string]interface{}{}
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Delete key-value pairs in map when the value is zero
	// for k, v := range incomingBookAsMap {
    //     if reflect.ValueOf(v).IsZero() {
    //         delete(incomingBookAsMap, k)
    //     }
    // }

	currentState := book.State // this is a pointer

	if incomingState, hasState := incomingBookAsMap["state"]; hasState {

		// Type assertion - needed because bookAsMap values are of type interface{}
		incomingState := incomingState.(string) // Type assertion
		incomingISBN := incomingBookAsMap["isbn"].(string) // Type assertion

		var incomingRequest *Book = &Book{&incomingISBN, ToPtr(incomingState), nil, nil, nil, nil}

		if incomingOnHoldCustomerID, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]; hasOnHoldCustomerID {
			incomingOnHoldCustomerID := incomingOnHoldCustomerID.(string) // Type assertion
			incomingRequest.OnHoldCustomerID = ToPtr(incomingOnHoldCustomerID)
		}

		if incomingCheckedOutCustomerID, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]; hasCheckedOutCustomerID {
			incomingCheckedOutCustomerID := incomingCheckedOutCustomerID.(string) // Type assertion
			incomingRequest.CheckedOutCustomerID = ToPtr(incomingCheckedOutCustomerID)
		}

		book, err = actionTable[*currentState][incomingState](book, incomingRequest)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
			return
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "MISSING STATE IN THE INCOMING REQUEST"})
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
		// curl -X PATCH localhost:8080/books/00 -H 'Content-Type: application/json' -H 'Accept: application/json' -d @incomingRequest.json