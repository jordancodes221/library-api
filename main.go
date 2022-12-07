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

type book struct{
	ISBN 					string 	`json:"isbn"`
	State 					string 	`json:"state"`

	Onhold_customerID 		string 	`json:"onhold_customerid"`
	Checkedout_customerID 	string `json:"checkedout_customerid"`

	Time_created 			time.Time `json:"time_created"`
	Time_updated  			time.Time `json:"time_updated"`
}

// Test data
var book_0 book = book{"0000", "available", 	"", 	"", 	time.Now(), time.Time{}} // not on-hold, not checked-out
var book_1 book = book{"0001", "checked-out", 	"", 	"01", 	time.Now(), time.Time{}} // checked-out, not on-hold
var book_2 book = book{"0002", "checked-out", 	"02", 	"01", 	time.Now(), time.Time{}} // checked-out by one customer, on-hold by another

var map_of_books = map[string]*book{
	"0000" : &book_0,
	"0001" : &book_1,
	"0002" : &book_2,
}

// GET (all books)
func getAllBooks(c *gin.Context) {
	// Make a slice containing all the values from map_of_books
	var vals []*book
	
	for _, v := range map_of_books {
		vals = append(vals, v)
	}

	c.IndentedJSON(http.StatusOK, vals)
}

// Helper function for GET (individual book)
func bookByISBN(isbn string) (*book, error) {
	book_ptr, ok := map_of_books[isbn]

	if ok {
		return book_ptr, nil
	} else {
		return nil, errors.New("Book not found.")
	}
}

// GET (individual book)
func getIndividualBook(c *gin.Context) {
	isbn := c.Param("isbn")
	book, err := bookByISBN(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// The following struct is needed to handle POST requests
type book_input struct{
	ISBN 					string 	`json:"isbn"`
	State 					string 	`json:"state"`
}

// POST
func createBook(c *gin.Context) {
	var newBookInput book_input

	if err := c.BindJSON(&newBookInput); err != nil {
		return // BindJSON handles the error response
	}

	var new_isbn string = newBookInput.ISBN
	var new_state string = newBookInput.State
	var newBook book = book{new_isbn, new_state, "", "", time.Now(), time.Time{}}
	map_of_books[new_isbn] = &newBook

	c.IndentedJSON(http.StatusCreated, newBook)
}

// DELETE
func deleteBook(c* gin.Context) {
	isbn := c.Param("isbn")
	delete(map_of_books, isbn)

	// The following prints the newly-shortened list of books
	// It makes this function O(N) rather than O(1), so perhaps it should be omitted
	var vals []*book
	for _, v := range map_of_books {
		vals = append(vals, v)
	}
	// End of possibly un-needed section

	c.IndentedJSON(http.StatusOK, vals)
}

// The following struct is needed to handle PATCH requests
type request struct{
	Requested_State string 		`json:"requested_state"`
	CustomerID 		string 		`json:"customerid"`

}

// PATCH
func updateBook(c* gin.Context) {
	isbn := c.Param("isbn")

	book, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Book not found."})
		return
	}
	current_state := book.State
	book_onhold_customerID := book.Onhold_customerID
	book_checkedout_customerID := book.Checkedout_customerID

	var newRequest request
	if err := c.BindJSON(&newRequest); err != nil {
		return
	}	
	requested_state := newRequest.Requested_State
	request_customerID := newRequest.CustomerID

	// Ensure requested state is valid
	if (requested_state != "available") && (requested_state != "on-hold") && (requested_state != "checked-out") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Invalid state requested."})
		return
	}

	if (current_state == "available") {
		book.State = requested_state
		if (requested_state == "available") {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "This book is already available."})
		}
		if (requested_state == "on-hold") {
			book.Onhold_customerID = request_customerID
			book.Time_updated = time.Now()
		}
		if (requested_state == "checked-out") {
			book.Checkedout_customerID = request_customerID
			book.Time_updated = time.Now()
		}
	}
	
	if (current_state == "on-hold") {
		if (book_onhold_customerID == request_customerID) { // The request comes from the customer who has it on-hold.
			book.State = requested_state
			book.Time_updated = time.Now()
			if (requested_state == "available") {
				book.Onhold_customerID = ""
				book.Time_updated = time.Now()
			}
			if (requested_state == "on-hold") {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "This customer already has this book on-hold."})
			}
			if (requested_state == "checked-out") {
				book.Onhold_customerID = ""
				book.Checkedout_customerID = request_customerID
				book.Time_updated = time.Now()
			}
		} else { // The request comes from a customer different from the one who has the book on-hold.
			// If the book's state is on-hold, no customer (other than the one who has it on-hold) can change its state
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book on-hold."})
		}
	} 
	
	if (current_state == "checked-out") {
		if (book_checkedout_customerID == request_customerID) { // The request comes from the customer who has it checked-out
			if (requested_state == "available") {
				// In this case, the customer who has it checked-out wishes to return the book (requesting the state is changed to "available")
				// If another customer has the book on-hold, the book's state will not be changed to "available" but rather to "on-hold" for that customer.
				if (book_onhold_customerID == "") { // Nobody has the book on-hold
					book.State = "available"
					book.Checkedout_customerID = ""
					book.Time_updated = time.Now()
				} else { // Another customer has the book on-hold
					book.State = "on-hold"
					book.Checkedout_customerID = ""
					book.Time_updated = time.Now()
				}
			}
			if (requested_state == "on-hold") {
				// In this case, the customer who has it checked-out wishes to place it on-hold.
				// If another customer has the book on-hold, this state change cannot be done.
				if (book_onhold_customerID == "") { // Nobody has the book on-hold
					book.State = "on-hold"
					book.Onhold_customerID = request_customerID
					book.Checkedout_customerID = ""
					book.Time_updated = time.Now()
				} else { // Another customer has the book on-hold
					c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book on-hold."})
				}
			}
			if (requested_state == "checked-out") {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "This customer already has this book checked-out."})
			}
		} else { // The request comes from a customer different from the one who has it checked-out
			if (requested_state == "available") {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book checked-out."})
			}
			if (requested_state == "on-hold") {
				// In this case, the book is checked-out and another customer requests to put it on-hold (meanwhile no other customer has it on-hold)
				// We can update the Onhold_customerID to accomodate this request, but the book's state will remain unchanged (it stays as "checked-out")
				if (book_onhold_customerID == "") { // Nobody has the book on-hold
					book.Onhold_customerID = request_customerID
					book.Time_updated = time.Now()
				} else { // Another customer has the book on-hold
					c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. The book is checked-out, and another customer has it on-hold."})
				}
			}
			if (requested_state == "checked-out") {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book checked-out. Consider trying to place it on-hold"})
			}
		}
	}
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()
	router.GET("/books", getAllBooks)
	router.GET("/books/:isbn", getIndividualBook)
	router.POST("/books", createBook)
	router.DELETE("/books/:isbn", deleteBook)
	router.PATCH("/books/:isbn", updateBook)

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
		// curl -X PATCH localhost:8080/books/0000 -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"Requested_State": "checked-out", "CustomerID": "01"}'