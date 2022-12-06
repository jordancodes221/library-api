// TO DO
	// add created and updated timestamps

	package main

	import (
		"net/http"
		"github.com/gin-gonic/gin"
		"errors"
		"time"
		// "encoding/json"
		//"strconv"
		// "fmt"
	)
	
	type book struct{
		ISBN 					string 	`json:"isbn"`
		State 					string 	`json:"state"` // available, checked-out

		Onhold_customerID 		string 	`json:"onhold_customerid"` // if nobody, "00"
		Checkedout_customerID 	string `json:"checkedout_customerid"`

		Time_created 			time.Time `json:"time_created"`
		Time_updated  			time.Time `json:"time_updated"`
	}

	type book_input struct{
		ISBN 					string 	`json:"isbn"`
		State 					string 	`json:"state"` // available, checked-out
	}
	
	// change books to be stored in a map, rather than an array
	var book_0 book = book{"0000", "available", 	"", 	"", 	time.Now(), time.Time{}} // not on-hold, not checked-out
	var book_1 book = book{"0001", "checked-out", 	"", 	"01", 	time.Now(), time.Time{}} // checked-out, not on-hold
	var book_2 book = book{"0002", "checked-out", 	"02", 	"01", 	time.Now(), time.Time{}} // checked-out, on-hold by another customer
	
	var map_of_books = map[string]*book{
		"0000" : &book_0,
		"0001" : &book_1,
		"0002" : &book_2,
	}
	
	// GET ALL
	func getAllBooks(c *gin.Context) {
		// get a slice of the all the values from the m map
		var vals []*book
		
		for _, v := range map_of_books {
			vals = append(vals, v) // *book and v... or book and *v???
		}
	
		c.IndentedJSON(http.StatusOK, vals)
	}
	
	// helper function -- changed this so it accesses map by key, and returns error if key not found
	func bookByISBN(isbn string) (*book, error) {
		book_ptr, ok := map_of_books[isbn]
	
		if ok {
			return book_ptr, nil
		} else {
			return nil, errors.New("Book not found.")
		}
	}
	
	// GET INDIVIDUAL
	func getIndividualBook(c *gin.Context) {
		isbn := c.Param("isbn")
		book, err := bookByISBN(isbn)
	
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Book not found."})
			return
		}
	
		c.IndentedJSON(http.StatusOK, book)
	}
	
	// POST
	func createBook(c *gin.Context) {
		var newBookInput book_input
	
		if err := c.BindJSON(&newBookInput); err != nil {
			return // BindJSON will handle the error response
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
	
		// probably don't actually need this part here  - don't need to re-print everything
		// want to get rid of it b/c it makes this delete function O(N) instead of O(1)
		var vals []*book
		for _, v := range map_of_books {
			vals = append(vals, v)
		}
		// end of possibly un-needed section
	
		c.IndentedJSON(http.StatusOK, vals)
	}
	
	// update this struct to be named "request"
	// make sure it includes customer ID
	type request struct{
		Requested_State string 	`json:"requested_state"`
		// request_time - date/time
		CustomerID 		string 	`json:"customerid"` // if nobody, "00"
	}
	
	// PATCH
	// CHANGE new_state TO requested_state
	func updateBook(c* gin.Context) {
		isbn := c.Param("isbn")
	
		book, err := bookByISBN(isbn)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Book not found."})
			return
		}
		current_state := book.State // maybe we should use a pointer instead???
		book_onhold_customerID := book.Onhold_customerID
		book_checkedout_customerID := book.Checkedout_customerID
	
		var newRequest request
		if err := c.BindJSON(&newRequest); err != nil {
			// probably also need some error checking so requested state is valid
			return // BindJSON will handle the error response
		}	
		new_state := newRequest.Requested_State
		request_customerID := newRequest.CustomerID
	
		// Business logic --> NEED TO REVISE THIS PART
		if (current_state == "available") {
			book.State = new_state
			if (new_state == "on-hold") {
				book.Onhold_customerID = request_customerID
			} else if (new_state == "checked-out") {
				book.Checkedout_customerID = request_customerID
			}
	
		} else if (current_state == "on-hold") {
			if (book_onhold_customerID == request_customerID) {
				// the request comes from the customer who has it on-hold
				// this customer can check it out or return it - both cases are handled below
				// the customer could also do an idempotent operation of putting it on-hold
				book.State = new_state
				// update the "updated time" field of book
				if (new_state == "checked-out") {
					book.Onhold_customerID = ""
					book.Checkedout_customerID = request_customerID
				}
				if (new_state == "available") {
					book.Onhold_customerID = ""
				}
			} else { // (book_onhold_customerID != request_customerID) && (current_state == "on-hold")
				// if a customer has the book on-hold, only that customer can check it out or return it
				// we handled both of these cases above
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book on-hold."})
			}
		} else { // current_state == "checked-out"}
			if (book_checkedout_customerID == request_customerID) {
				// the book is currently checked out, and a request comes in from the customer who has it
				// this customer can either return it, or place it on-hold if nobody else
				// this case is more complicated...
					// the state cannot be determined simply because the ID's match, unlike the on-hold case
					// lets say the requested state is available; if someone has it on-hold then it goes to OH
				if (new_state == "available") {
					if (book_onhold_customerID == "") { // nobody has it on-hold
						book.State = "available"
						book.Checkedout_customerID = ""
					} else { // somebody does have it on-hold
						book.State = "on-hold"
						book.Checkedout_customerID = ""
					}
				} else if (new_state == "on-hold") {
					if (book_onhold_customerID == "") { // nobody has it on-hold
						book.State = "on-hold"
						book.Onhold_customerID = request_customerID
						book.Checkedout_customerID = ""
					} else { // somebody does have it on-hold
						c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book on-hold."})
					}
				} else { // new state is checked-out and customers id
					book.State = "checked-out" // don't really need this, but didn't want to leave it empty
				}
			} else { // book is checked-out, but somebody else makes a request
				// if requested state is on-hold but book is checked out, we cannot change the state to
				// on-hold, but what we can do is update the on-hold customer ID if it is currently ""
				// i.e. if nobody has it on-hold, the customer can put it on hold without changing its state
				if (new_state == "available") {
					c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book checked-out."})
				} else if (new_state == "on-hold") {
					// this is the most complex case
					if (book_onhold_customerID == "") { // nobody has it on-hold
						book.Onhold_customerID = request_customerID
					} else { // somebody does have it on-hold
						c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. The book is checked-out, and another customer has it on-hold."})
					}
				} else { // requested state is to check it out... but it's already checked out
					c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Request failed. Another customer has the book checked-out. Consider placing it on-hold"})
				}
			}
		}
	
		// how to handle invalid state --> may not have to if you we give drop-down menu on front end??
	
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
	
	// PATCH request is as follows
	// had an error because the request had State instead of Requested_State in the json
	// NEED TO UPDATE PATCH REQUEST DATA TO MATCH NEW SCHEMA .... NEVERMIND!!!!!
		// curl -X PATCH localhost:8080/books/0000 -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"Requested_State": "checked-out", "CustomerID": "01"}'
	
	// POST
		// curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"
	////////////////////////////////////////////////////////////////////////////////////