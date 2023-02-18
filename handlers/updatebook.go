package handlers

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
	if err := models.ValidateISBNAndStateSyntax(incomingBookAsMap); err != nil {
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
