package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	"encoding/json"
	"fmt"
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"details": "REQUEST SUCCESSFUL. BOOK NOT FOUND"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// Syntactic Validation
func ValidateISBNAndStateSyntax(incomingBookAsMap map[string]interface{}) (error) { // the request will not complete if input is not OK, which why it is possible to return an error		
	// Assuming ISBN is present, is it valid?
	if isbn, hasISBN := incomingBookAsMap["isbn"]; hasISBN {
		_, isbnIsString := isbn.(string)
		if !isbnIsString {
			return errors.New("ISBN provided is not of type string.")
		}
	}

	// Assuming State is present, is it valid?
	if state, hasState := incomingBookAsMap["state"]; hasState {
		state, stateIsString := state.(string)
		if !stateIsString {
			return errors.New("State provided is not of type string.")
		}

		if ((state != "available") && (state != "on-hold") && (state != "checked-out")) {
			return errors.New("Invalid state provided. State must be equal to one of: \"available\", \"on-hold\", or \"checked-out\".")
		}
	}

	return nil
}

func ValidateIDSemanticsForCreateBook(incomingBookAsMap map[string]interface{}) (error) {
	fmt.Println("CALLING ValidateIDSemanticsForCreateBook...")

	// This function will only be called once state is established to be both present and valid
	state := incomingBookAsMap["state"]
	
	// Retrieve the customer ID's if they are present
	_, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]
	_, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]

	// State is available -- THIS IS SEMANTIC CHECKING
	if (state == "available") {
		if ((hasOnHoldCustomerID) && (!hasCheckedOutCustomerID)) {
			return errors.New("Cannot have an on-hold customer ID when state is available.")
		}

		if (!(hasOnHoldCustomerID) && (hasCheckedOutCustomerID)) {
			return errors.New("Cannot have checked-out customer ID when state is available.")
		}
		
		if (hasOnHoldCustomerID || hasCheckedOutCustomerID) {
			return errors.New("Cannot have on-hold customer ID or checked-out customer ID when state is available.")
		}
	}

	// State is on-hold -- THIS IS SEMANTIC CHECKING
	if (state == "on-hold") {
		if hasCheckedOutCustomerID {
			return errors.New("Cannot have checked-out customer ID when state is on-hold.")
		}

		if hasOnHoldCustomerID {
			// We know ohid is provided. Ensure it is a string
			ohid, ohidIsString := incomingBookAsMap["onholdcustomerid"].(string)
			if !ohidIsString {
				return errors.New("On-hold customer ID provided is not of type string.")
			}

			if (ohid == "") {
				return errors.New("On-hold customer ID is the empty string.")
			}
		} else { // !hasOnHoldCustomerID
			return errors.New("State provided is on-hold, but no on-hold customer ID is provided.")
		}
	}

	// State is checked-out -- THIS IS SEMANTIC CHECKING
	if (state == "checked-out") {
		if hasOnHoldCustomerID {
			return errors.New("Cannot have on-hold customer ID when state is checked-out.")
		}

		if hasCheckedOutCustomerID {
			// We know ohid is provided. Ensure it is a string
			coid, coidIsString := incomingBookAsMap["checkedoutcustomerid"].(string)
			if !coidIsString {
				return errors.New("Checked-out customer ID provided is not of type string.")
			}

			if (coid == "") {
				return errors.New("Checked-out customer ID is the empty string.")
			}
		} else { // !hasCheckedOutCustomerID
			return errors.New("State provided is checked-out, but no checked-out customer ID is provided.")
		}
	}

	return nil
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

	// Ensure that incoming JSON includes ISBN
	if _, hasISBN := incomingBookAsMap["isbn"]; !hasISBN {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Missing ISBN in the incoming request."})
		return
	}

	// The reason for calling validate at this point is that it is:
		// (1) After checking that ISBN is present, and
		// (2) Before checking if ISBN is in-use (we want to ensure it's valid before checking if it's in-use)

	// Validate ISBN
	if err := ValidateISBNAndStateSyntax(incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Make sure ISBN is not already in-use
		// At this point, we know that ISBN (1) is present, and (2) is valid
	incomingISBN := incomingBookAsMap["isbn"].(string)
	if _, ok := mapOfBooks[incomingISBN]; ok {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": "Book already exists."})
		return
	}

	// Update newBook ISBN field
	newBook.ISBN = ToPtr(incomingISBN)

	// Update newBook State field (if state is present)
	if incomingState, hasState := incomingBookAsMap["state"]; hasState {
		incomingState := incomingState.(string) // Type Assertion
		newBook.State = ToPtr(incomingState)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Missing State in the incoming request."})
		return
	}

	// Ensure correct customer ID fields are provided given the state
	if err := ValidateIDSemanticsForCreateBook(incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Update newBook OnHoldCustomerID
	if incomingOnHoldCustomerID, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]; hasOnHoldCustomerID {
		incomingOnHoldCustomerID := incomingOnHoldCustomerID.(string) // Type assertion
		newBook.OnHoldCustomerID = ToPtr(incomingOnHoldCustomerID)
	}

	// Update newBook CheckedOutCustomerID
	if incomingCheckedOutCustomerID, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]; hasCheckedOutCustomerID {
		incomingCheckedOutCustomerID := incomingCheckedOutCustomerID.(string) // Type assertion
		newBook.CheckedOutCustomerID = ToPtr(incomingCheckedOutCustomerID)
	}

	// Update newBook times
	newBook.TimeCreated = ToPtr(time.Now())
	newBook.TimeUpdated = nil

	// Add newBook to mapOfBooks
	mapOfBooks[*newBook.ISBN] = newBook

	c.IndentedJSON(http.StatusCreated, newBook) // 201 status code if successful
}

// DELETE
func DeleteBook(c *gin.Context) {
	isbn := c.Param("isbn")

	book, err := bookByISBN(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code
		return
	}

	if book == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"details": "REQUEST SUCCESSFUL. BOOK NOT FOUND"})
		return
	}

	delete(mapOfBooks, isbn)
	c.Status(http.StatusNoContent) // 204 status code if successful
}

// Semantic Validation for Checkout and Return
func ValidateIDSemanticsForCheckedOutUpdate(incomingRequest *Book) (error) {
	// incomingRequest is of the form &{isbn, state, checkedoutcustomerid, onholdcustomerid, timecreated, timeupdated}
	// For this particular case, it should be populated as such: &{isbn, state, checkedoutcustomerid, nil, nil, nil}
	
	checkedOutCustomerID := incomingRequest.CheckedOutCustomerID
	onHoldCustomerID := incomingRequest.OnHoldCustomerID

	if (checkedOutCustomerID == nil && onHoldCustomerID == nil) {
		return errors.New("Expected checked-out customer ID.")
	}

	if (checkedOutCustomerID != nil && onHoldCustomerID != nil) {
		return errors.New("Did not expected on-hold customer ID.")
	}

	if (checkedOutCustomerID == nil && onHoldCustomerID != nil) {
		return errors.New("Expected checked-out customer ID, and did not expected on-hold customer ID.")
	}

	return nil
}

// Semantic Validation for PlaceHold and ReleaseHold
func ValidateIDSemanticsForOnHoldUpdate(incomingRequest *Book) (error) {
	// incomingRequest is of the form &{isbn, state, checkedoutcustomerid, onholdcustomerid, timecreated, timeupdated}
	// For this particular case, it should be populated as such: &{isbn, state, nil, onholdcustomerid, nil, nil}
	checkedOutCustomerID := incomingRequest.CheckedOutCustomerID
	onHoldCustomerID := incomingRequest.OnHoldCustomerID

	if (onHoldCustomerID == nil && checkedOutCustomerID == nil) {
		return errors.New("Expected on-hold customer ID.")
	}

	if (onHoldCustomerID != nil && checkedOutCustomerID != nil) {
		return errors.New("Did not expected checked-out customer ID.")
	}

	if (onHoldCustomerID == nil && checkedOutCustomerID != nil) {
		return errors.New("Expected on-hold customer ID, and did not expected checked-out customer ID.")
	}

	return nil
}

// No Match Error
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

	book, err := bookByISBN(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if book == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"details": "REQUEST SUCCESSFUL. BOOK NOT FOUND"})
		return
	}	

	incomingBookAsMap := map[string]interface{}{}
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Missing state in the incoming request."})
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