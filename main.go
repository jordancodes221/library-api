package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	"fmt"
	"encoding/json"
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

// Test data
func ToPtr[T string|time.Time](v T) *T {
    return &v
}

var bookInstance00 Book = Book{ToPtr("00"), ToPtr("available"), nil, nil, ToPtr(time.Now()), ToPtr(time.Now())}

var mapOfBooks = map[string]*Book{
	"00" : &bookInstance00,
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

// POST
func CreateBook(c *gin.Context) {
	var newBook Book // value type... this is allocating memory for the book - not a terrible thing
	// but typically when working with structs, you want to declare the struct variables to be pointer types
	// the reason you make it a pointer type is because there are 2 different ways values get passed to functions:
		// 1) pass by value, 
		// 2) pass by reference
	// in general, you never want to pass a struct by value, because you don't want the overhead of copying the struct
	// the other thing that can happen is the thing re-assigning hte pointers when making copies...

	// Unmarshal
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Make sure ISBN is not already in-use
	if _, ok := mapOfBooks[*newBook.ISBN]; ok {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": "Book already exists."})
		return
	}

	*newBook.TimeCreated = time.Now()
	*newBook.TimeUpdated = time.Now()

	mapOfBooks[*newBook.ISBN] = &newBook

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
	fmt.Println("CHECKING OUT...")
	if (*currentBook.State == "available") {
		fmt.Println("CHECKING OUT... CURRENT STATE IS AVAILABLE")
		*currentBook.State = "checked-out" // or should we use incomingBook.State? 
		fmt.Println("FIRST LINE COMPLETE")
		currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
		fmt.Println("SECOND LINE COMPLETE")
		newTimeUpdated := time.Now()
		*currentBook.TimeUpdated = newTimeUpdated
		fmt.Println("THIRD LINE COMPLETE")
	}

	if (*currentBook.State == "on-hold") {
		if (currentBook.OnHoldCustomerID == incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			*currentBook.State = "checked-out"
			currentBook.OnHoldCustomerID = nil
			currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
			*currentBook.TimeUpdated = time.Now()
		} else {
			return nil, NoMatchError
		}
	}

	if (*currentBook.State == "checked-out") {
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
	if (*currentBook.State == "available") {
		*currentBook.State = "on-hold"
		currentBook.OnHoldCustomerID = incomingBook.OnHoldCustomerID
		*currentBook.TimeUpdated = time.Now()
	}

	if (*currentBook.State == "on-hold") {
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
	if (*currentBook.State == "on-hold") {
		if (currentBook.OnHoldCustomerID == incomingBook.OnHoldCustomerID) {
			*currentBook.State = "available"
			currentBook.OnHoldCustomerID = nil
			*currentBook.TimeUpdated = time.Now()
		} else {
			return nil, NoMatchError
		}
	}

	return currentBook, nil
}

// Return
	// checked-out --> available (when ID's match)
func Return(currentBook *Book, incomingBook *Book) (*Book, error) {
	if (*currentBook.State == "checked-out") {
		if (currentBook.CheckedOutCustomerID == incomingBook.CheckedOutCustomerID) {
			*currentBook.State = "available"
			currentBook.CheckedOutCustomerID = nil // need this, or leave it as who most recently had it on hold?
			*currentBook.TimeUpdated = time.Now()
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
	fmt.Println("\nUPDATE BOOK CALLED")
	isbn := c.Param("isbn")

	// Ensure book to be updated exists
	book, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": err.Error()})
		return
	}

	fmt.Println("BEGINNING BOOK AS MAP DECODING")

	incomingBookAsMap := map[string]interface{}{}
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	fmt.Println("ENDING BOOK AS MAP DECODING")

	fmt.Println("\nPRINTING BOOK AS MAP...")
	fmt.Println("")
	fmt.Println(incomingBookAsMap)
	fmt.Println("")
	fmt.Println("\nDONE PRINTING BOOK AS MAP.")


	currentState := book.State // this is is a pointer

	if _, hasState := incomingBookAsMap["state"]; hasState {

		// Type assertion
		incomingState := incomingBookAsMap["state"].(string)
		incomingISBN := incomingBookAsMap["isbn"].(string)

		var incomingRequest Book = Book{&incomingISBN, &incomingState, nil, nil, nil, nil}

		var incomingCheckedOutCustomerIDptr *string
		if _, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]; hasCheckedOutCustomerID {
			incomingCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"].(string)
			incomingCheckedOutCustomerIDptr = &incomingCheckedOutCustomerID
		}
		incomingRequest.CheckedOutCustomerID = incomingCheckedOutCustomerIDptr

		var incomingOnHoldCustomerIDptr *string
		if _, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]; hasOnHoldCustomerID {
			incomingOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"].(string)
			incomingRequest.OnHoldCustomerID = &incomingOnHoldCustomerID
		}
		incomingRequest.OnHoldCustomerID = incomingOnHoldCustomerIDptr

		book, err = actionTable[*currentState][incomingState](book, &incomingRequest)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
			return
		}
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