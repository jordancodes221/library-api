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

	c.IndentedJSON(http.StatusCreated, newBook)
}

// DELETE
func DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")
	delete(mapOfBooks, isbn)
	c.Status(http.StatusNoContent) // 204 status code
}

// The following struct is needed to handle PATCH requests
type Request struct{
	RequestedState string 		`json:"requestedstate"`
	CustomerID 		string 		`json:"customerid"`

}

// PATCH
func UpdateBook(c *gin.Context) {
	isbn := c.Param("isbn")

	book, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Book not found."})
		return
	}
	currentState := book.State

	var newRequest Request
	if err := c.BindJSON(&newRequest); err != nil {
		return
	}	
	requestedState := newRequest.RequestedState
	requestCustomerID := newRequest.CustomerID

	// Ensure requested state is valid
	if (requestedState != "available") && (requestedState != "on-hold") && (requestedState != "checked-out") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Invalid state requested."})
		return
	}

	if (currentState == "available") {
		book.State = requestedState
		if (requestedState == "available") {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "This book is already available."})
		}
		if (requestedState == "on-hold") {
			book.OnHoldCustomerID = requestCustomerID
			book.TimeUpdated = time.Now().String()
		}
		if (requestedState == "checked-out") {
			book.CheckedOutCustomerID = requestCustomerID
			book.TimeUpdated = time.Now().String()
		}
	}
	
	if (currentState == "on-hold") {
		if (book.OnHoldCustomerID == requestCustomerID) { // The request comes from the customer who has it on-hold.
			book.State = requestedState
			book.TimeUpdated = time.Now().String()
			if (requestedState == "available") {
				book.OnHoldCustomerID = ""
				book.TimeUpdated = time.Now().String()
			}
			if (requestedState == "on-hold") {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "This customer already has this book on-hold."})
			}
			if (requestedState == "checked-out") {
				book.OnHoldCustomerID = ""
				book.CheckedOutCustomerID = requestCustomerID
				book.TimeUpdated = time.Now().String()
			}
		} else { // The request comes from a customer different from the one who has the book on-hold.
			// If the book's state is on-hold, no customer (other than the one who has it on-hold) can change its state
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book on-hold."})
		}
	} 
	
	if (currentState == "checked-out") {
		if (book.CheckedOutCustomerID == requestCustomerID) { // The request comes from the customer who has it checked-out
			if (requestedState == "available") {
				// In this case, the customer who has it checked-out wishes to return the book (requesting the state is changed to "available")
				// If another customer has the book on-hold, the book's state will not be changed to "available" but rather to "on-hold" for that customer.
				if (book.OnHoldCustomerID == "") { // Nobody has the book on-hold
					book.State = "available"
					book.CheckedOutCustomerID = ""
					book.TimeUpdated = time.Now().String()
				} else { // Another customer has the book on-hold
					book.State = "on-hold"
					book.CheckedOutCustomerID = ""
					book.TimeUpdated = time.Now().String()
				}
			}
			if (requestedState == "on-hold") {
				// In this case, the customer who has it checked-out wishes to place it on-hold.
				// If another customer has the book on-hold, this state change cannot be done.
				if (book.OnHoldCustomerID == "") { // Nobody has the book on-hold
					book.State = "on-hold"
					book.OnHoldCustomerID = requestCustomerID
					book.CheckedOutCustomerID = ""
					book.TimeUpdated = time.Now().String()
				} else { // Another customer has the book on-hold
					c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book on-hold."})
				}
			}
			if (requestedState == "checked-out") {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "This customer already has this book checked-out."})
			}
		} else { // The request comes from a customer different from the one who has it checked-out
			if (requestedState == "available") {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book checked-out."})
			}
			if (requestedState == "on-hold") {
				// In this case, the book is checked-out and another customer requests to put it on-hold (meanwhile no other customer has it on-hold)
				// We can update the OnHoldCustomerID to accomodate this request, but the book's state will remain unchanged (it stays as "checked-out")
				if (book.OnHoldCustomerID == "") { // Nobody has the book on-hold
					book.OnHoldCustomerID = requestCustomerID
					book.TimeUpdated = time.Now().String()
				} else { // Another customer has the book on-hold
					c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. The book is checked-out, and another customer has it on-hold."})
				}
			}
			if (requestedState == "checked-out") {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book checked-out. Consider trying to place it on-hold"})
			}
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
		// curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"
	// DELETE
		// curl localhost:8080/books/0000 --request "DELETE"
	// PATCH
		// curl -X PATCH localhost:8080/books/0000 -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"RequestedState": "checked-out", "CustomerID": "01"}'