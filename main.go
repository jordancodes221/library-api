// TO DO
	// change schema so theres 2 customerID's - one for on-hold, one for checked-out
	// need to change the business logic a little bit because of this...
		// if book is checked-out and is then returned, it will still have an on-hold customer
		// so book can only be checked-out by that customer
	// add created and updated timestamps also



	package main

	import (
		"net/http"
		"github.com/gin-gonic/gin"
		"errors"
		// "encoding/json"
		//"strconv"
		// "fmt"
	)
	
	type book struct{
		ISBN 		string 	`json:"isbn"`
		State 		string 	`json:"state"` // available, checked-out
		CustomerID 	string 	`json:"customerid"` // if "00", nobody checked it out
	}
	
	// change books to be stored in a map, rather than an array
	var book_0 book = book{"0000", "available", "00"}
	var book_1 book = book{"0001", "checked-out", "01"}
	var book_2 book = book{"0002", "on-hold", "02"}
	
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
		var newBook book
	
		if err := c.BindJSON(&newBook); err != nil {
			return // BindJSON will handle the error response
		}
	
		// GET ISBN NUMBER FROM NEWBOOK JSON then map_of_books[newbook_isbn] = newbook
		var new_isbn string = newBook.ISBN
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
		CustomerID 		string	`json:"customerid"`
	}
	
	// PATCH
	func updateBook(c* gin.Context) {
		isbn := c.Param("isbn")
	
		book, err := bookByISBN(isbn)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Book not found."})
			return
		}
		old_state := book.State // maybe we should use a pointer instead???
		book_customerID := book.CustomerID
	
		var newRequest request
		if err := c.BindJSON(&newRequest); err != nil {
			// probably also need some error checking so requested state is valid
			return // BindJSON will handle the error response
		}	
		new_state := newRequest.Requested_State
		request_customerID := newRequest.CustomerID
	
		// Business logic
		if (old_state == "available") {
			book.State = new_state
			if (new_state != "available") {
				book.CustomerID = request_customerID
				// if you update available --> available, you can't manipulate customer id
			}
		} else { // if (old_state == "on-hold") || (old_state == "checked-out")
		///// maybe ON-HOLD is if you request to check-out an already checked out book?
			if (book_customerID == request_customerID) {
				book.State = new_state
				if (new_state == "available") {
					book.CustomerID = "00"
				}
			} else { // book_customerID != request_customerID
				InvalidStateTransfer_message := "Cannot update state to " + new_state + " because another customer already has it " + old_state + "."
				c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": InvalidStateTransfer_message})
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
		// curl -X PATCH localhost:8080/books/0000 -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"Requested_State": "checked-out", "CustomerID": "01"}'
	
	////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////
	
	// // CHECKOUT BOOK
	// func checkoutBook(c *gin.Context) {
	// 	isbn, ok := c.GetQuery("isbn")
	
	// 	if ok == false {
	// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Query parameter (ISBN) is missing."})
	// 		return
	// 	}
	
	// 	book, err := bookByISBN(isbn)
	// 	if err != nil {
	// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
	// 		return
	// 	}
	// 	if book.State == "checked-out" {
	// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book is already checked-out."})
	// 		return
	// 	}
	
	// 	book.State = "checked-out"
	// 	c.IndentedJSON(http.StatusOK, book)
	// }
	
	// RETURN BOOK
	// func returnBook(c *gin.Context) {
	// 	isbn, ok := c.GetQuery("isbn")
	// 	if ok == false {
	// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Query parameter (ISBN) is missing."})
	// 		return
	// 	}
	
	// 	book, err := bookByISBN(isbn)
	// 	if err != nil {
	// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
	// 		return
	// 	}
	// 	if book.State == "available" {
	// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book has already been returned."})
	// 		return
	// 	}
	
	// 	book.State = "available"
	// 	c.IndentedJSON(http.StatusOK, book)
	// }
	
	
	