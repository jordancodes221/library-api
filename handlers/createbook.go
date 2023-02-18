package handlers

import ( 
	// "example/library_project/handlers"
	"example/library_project/validators"
	"example/library_project/models"
	


	"net/http"
	"github.com/gin-gonic/gin"
	// "errors"
	"time"
	"encoding/json"
	// "fmt"
	// "reflect"
	// "strconv"
)

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
	if err := models.ValidateISBNAndStateSyntax(incomingBookAsMap); err != nil {
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