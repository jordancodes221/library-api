package main

import ( 
	// "example/library_project/handlers"
	"example/library_project/validators"
	"example/library_project/models"
	


	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	"encoding/json"
	"fmt"
	// "reflect"
	// "strconv"
)

//// Instantiating Test Data
// Generic function converts literals to pointers
func ToPtr[T string|time.Time](v T) *T {
    return &v
}

// First test of instantiating test data with new schema and ToPtr function
var bookInstance00 models.Book = models.Book{ISBN: ToPtr("00"), State: ToPtr("on-hold"), OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Now())}

// Actual test data to be used in testing
var bookInstance0 models.Book = models.Book{ISBN: ToPtr("0000"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available
var bookInstance1 models.Book = models.Book{ISBN: ToPtr("0001"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out
var bookInstance2 models.Book = models.Book{ISBN: ToPtr("0002"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold

var bookInstance3 models.Book = models.Book{ISBN: ToPtr("0003"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available
var bookInstance4 models.Book = models.Book{ISBN: ToPtr("0004"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available (no match)
var bookInstance5 models.Book = models.Book{ISBN: ToPtr("0005"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out
var bookInstance6 models.Book = models.Book{ISBN: ToPtr("0006"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out (no match)
var bookInstance7 models.Book = models.Book{ISBN: ToPtr("0007"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold 
var bookInstance8 models.Book = models.Book{ISBN: ToPtr("0008"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold (no match)

var bookInstance9 models.Book =  models.Book{ISBN: ToPtr("0009"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available
var bookInstance10 models.Book = models.Book{ISBN: ToPtr("0010"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available (no match)
var bookInstance11 models.Book = models.Book{ISBN: ToPtr("0011"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out
var bookInstance12 models.Book = models.Book{ISBN: ToPtr("0012"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out (no match)
var bookInstance13 models.Book = models.Book{ISBN: ToPtr("0013"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold 
var bookInstance14 models.Book = models.Book{ISBN: ToPtr("0014"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold (no match)

var bookInstance15 models.Book = models.Book{ISBN: ToPtr("0015"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, 	TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> This is the book to be deleted in testing

// The following are for UpdateBook ID semantics validation
var bookInstance16 models.Book = models.Book{ISBN: ToPtr("0016"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, 	TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} 
var bookInstance17 models.Book = models.Book{ISBN: ToPtr("0017"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})}
var bookInstance18 models.Book = models.Book{ISBN: ToPtr("0018"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})}

// Map of test data to be used in testing
var mapOfBooks = map[string]*models.Book{
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

	"0016" : &bookInstance16,
	"0017" : &bookInstance17,
	"0018" : &bookInstance18,
}

// GET (all books)
func GetAllBooks(c *gin.Context) {
	// Make a slice containing all the values from mapOfBooks
	var vals []*models.Book
	
	for _, v := range mapOfBooks {
		vals = append(vals, v)
	}

	c.IndentedJSON(http.StatusOK, vals)
}

// Helper function for GET (individual book)
func bookByISBN(isbn string) (*models.Book, error) {
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

// POST
func CreateBook(c *gin.Context) {
	var newBook *models.Book = &models.Book{
		ISBN: nil, 
		State: nil, 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: nil, 
		TimeUpdated: nil,}

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

	// Validate ISBN and State Syntax
	if err := validators.ValidateISBNAndStateSyntax(incomingBookAsMap); err != nil {
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

	// Validate Time Semantics
	if err := validators.ValidateTimeSemantics(incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
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
	if err := validators.ValidateIDSemanticsForCreateBook(incomingBookAsMap); err != nil {
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

// No Match Error
var NoMatchError error = errors.New("ID's do not match.")

// Checkout
	// available --> checked-out
	// on-hold --> checked-out
	// checked-out --> checked-out
func Checkout(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := validators.ValidateIDSemanticsForCheckedOutUpdate(incomingBook); err != nil {
		return nil, err
	}

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
func Conflict(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	return nil, errors.New("Invalid state transfer requested.")
}

// PlaceHold
	// available --> on-hold
	// on-hold --> on-hold
func PlaceHold(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := validators.ValidateIDSemanticsForOnHoldUpdate(incomingBook); err != nil {
		return nil, err
	}

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
func ReleaseHold(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := validators.ValidateIDSemanticsForOnHoldUpdate(incomingBook); err != nil {
		return nil, err
	}

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
func Return(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	fmt.Println("CALLING RETURNBOOK...")

	if err := validators.ValidateIDSemanticsForCheckedOutUpdate(incomingBook); err != nil {
		return nil, err
	}

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
func NoOperation(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	return currentBook, nil
}

// First key is current state, 2nd key is incoming state
var actionTable = map[string]map[string]func(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
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

	currentBook, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if currentBook == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"details": "REQUEST SUCCESSFUL. BOOK NOT FOUND"})
		return
	}	

	incomingBookAsMap := map[string]interface{}{}
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Validate ISBN and State Syntax
	if err := validators.ValidateISBNAndStateSyntax(incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Validate Time Semantics
	if err := validators.ValidateTimeSemantics(incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	currentState := currentBook.State // this is a pointer

	if incomingState, hasState := incomingBookAsMap["state"]; hasState {

		// Type assertion - needed because currentBookAsMap values are of type interface{}
		incomingState := incomingState.(string) // Type assertion
		incomingISBN := incomingBookAsMap["isbn"].(string) // Type assertion

		// THE REASON WE HAVE A PROBLEM HERE IS THAT WE HAVE A LOCAL VARIABLE ALSO CALLED BOOK....
		var incomingRequest *models.Book = &models.Book{
			ISBN: &incomingISBN, 
			State: ToPtr(incomingState), 
			OnHoldCustomerID: nil, 
			CheckedOutCustomerID: nil, 
			TimeCreated: nil, 
			TimeUpdated: nil,
		}

		if incomingOnHoldCustomerID, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]; hasOnHoldCustomerID {
			incomingOnHoldCustomerID := incomingOnHoldCustomerID.(string) // Type assertion
			incomingRequest.OnHoldCustomerID = ToPtr(incomingOnHoldCustomerID)
		}

		if incomingCheckedOutCustomerID, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]; hasCheckedOutCustomerID {
			incomingCheckedOutCustomerID := incomingCheckedOutCustomerID.(string) // Type assertion
			incomingRequest.CheckedOutCustomerID = ToPtr(incomingCheckedOutCustomerID)
		}

		currentBook, err = actionTable[*currentState][incomingState](currentBook, incomingRequest)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
			return
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Missing state in the incoming request."})
		return
	}

	c.IndentedJSON(http.StatusOK, currentBook)
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