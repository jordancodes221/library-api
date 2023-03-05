package handlers

import ( // h.books, h.bookByISBN
	"example/library_project/models"
	"example/library_project/utils"
	
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	"encoding/json"
	"fmt"
)

// Validate Time Semantics
func ValidateTimeSemanticsForUpdateBook(currentBook *models.Book, incomingBookAsMap map[string]interface{}) (error) {
	fmt.Println("CALLING VALIDATE TIME SEMANTICS...")

	currentTimeCreated := *currentBook.TimeCreated
	if incomingTimeCreatedUnparsed, ok := incomingBookAsMap["timecreated"]; ok {
		incomingTimeCreatedUnparsed := incomingTimeCreatedUnparsed.(string)
		incomingTimeCreated, _ := time.Parse(time.RFC3339, incomingTimeCreatedUnparsed)
		
		fmt.Print("CURRENT TIME CREATED: ")
		fmt.Println(currentTimeCreated)
		fmt.Print("REQUESTED TIME CREATED: ")
		fmt.Println(incomingTimeCreated)

		if currentTimeCreated != incomingTimeCreated {
			return errors.New("Requested time created does not match existing time created.")
		}
	}

	currentTimeUpdated := *currentBook.TimeUpdated
	if incomingTimeUpdatedUnparsed, ok := incomingBookAsMap["timeupdated"]; ok {
		incomingTimeUpdatedUnparsed := incomingTimeUpdatedUnparsed.(string)
		incomingTimeUpdated, _ := time.Parse(time.RFC3339, incomingTimeUpdatedUnparsed)
		
		fmt.Print("CURRENT TIME UPDATED: ")
		fmt.Println(currentTimeUpdated)
		fmt.Print("REQUESTED TIME UPDATED: ")
		fmt.Println(incomingTimeUpdated)
		
		if currentTimeUpdated != incomingTimeUpdated {
			return errors.New("Requested time updated does not match existing time updated.")
		}
	}
	
	return nil
}


// Semantic Validation for Checkout and Return
func ValidateIDSemanticsForCheckedOutUpdate(incomingRequest *models.Book) (error) {
	// incomingRequest is of the form &{isbn, state, checkedoutcustomerid, onholdcustomerid, timecreated, timeupdated}
	// For this particular case, it should be populated as such: &{isbn, state, checkedoutcustomerid, nil, nil, nil}
	
	// fmt.Println("CALLING ValidateIDSemanticsForCheckedOutUpdate...")
	checkedOutCustomerID := incomingRequest.CheckedOutCustomerID
	onHoldCustomerID := incomingRequest.OnHoldCustomerID

	if (checkedOutCustomerID == nil) {
		return errors.New("Expected checked-out customer ID.")
	}

	if (onHoldCustomerID != nil) {
		return errors.New("Did not expect on-hold customer ID.")
	}

	return nil
}

// Semantic Validation for PlaceHold and ReleaseHold
func ValidateIDSemanticsForOnHoldUpdate(incomingRequest *models.Book) (error) {
	// incomingRequest is of the form &{isbn, state, checkedoutcustomerid, onholdcustomerid, timecreated, timeupdated}
	// For this particular case, it should be populated as such: &{isbn, state, nil, onholdcustomerid, nil, nil}
	checkedOutCustomerID := incomingRequest.CheckedOutCustomerID
	onHoldCustomerID := incomingRequest.OnHoldCustomerID

	// fmt.Println("CALLING ValidateIDSemanticsForOnHoldUpdate...")
	if (onHoldCustomerID == nil) {
		return errors.New("Expected on-hold customer ID.")
	}

	if (checkedOutCustomerID != nil) {
		return errors.New("Did not expect checked-out customer ID.")
	}

	return nil
}

// No Match Error
var NoMatchError error = errors.New("ID's do not match.")

// Checkout
	// available --> checked-out
	// on-hold --> checked-out
	// checked-out --> checked-out
func Checkout(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := ValidateIDSemanticsForCheckedOutUpdate(incomingBook); err != nil {
		return nil, err
	}

	if (*currentBook.State == "available") {
		*currentBook.State = "checked-out" // or should we use incomingBook.State? 
		currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
		currentBook.TimeUpdated = utils.ToPtr(time.Now())
	} else if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			*currentBook.State = "checked-out"
			currentBook.OnHoldCustomerID = nil
			currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
			currentBook.TimeUpdated = utils.ToPtr(time.Now())
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
	if err := ValidateIDSemanticsForOnHoldUpdate(incomingBook); err != nil {
		return nil, err
	}

	if (*currentBook.State == "available") {
		*currentBook.State = "on-hold"
		currentBook.OnHoldCustomerID = incomingBook.OnHoldCustomerID
		currentBook.TimeUpdated = utils.ToPtr(time.Now())
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
	if err := ValidateIDSemanticsForOnHoldUpdate(incomingBook); err != nil {
		return nil, err
	}

	if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.OnHoldCustomerID) {
			*currentBook.State = "available"
			currentBook.OnHoldCustomerID = nil
			currentBook.TimeUpdated = utils.ToPtr(time.Now())
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

	if err := ValidateIDSemanticsForCheckedOutUpdate(incomingBook); err != nil {
		return nil, err
	}

	if (*currentBook.State == "checked-out") {
		if (*currentBook.CheckedOutCustomerID == *incomingBook.CheckedOutCustomerID) {
			*currentBook.State = "available"
			currentBook.CheckedOutCustomerID = nil
			currentBook.TimeUpdated = utils.ToPtr(time.Now())
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
func (h *BooksHandler) UpdateBook(c *gin.Context) {
	isbn := c.Param("isbn")

	currentBook, err := h.bookByISBN(isbn)
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
	if err := ValidateTimeSemanticsForUpdateBook(currentBook, incomingBookAsMap); err != nil {
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
			State: utils.ToPtr(incomingState), 
			OnHoldCustomerID: nil, 
			CheckedOutCustomerID: nil, 
			TimeCreated: nil, 
			TimeUpdated: nil,
		}

		if incomingOnHoldCustomerID, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]; hasOnHoldCustomerID {
			incomingOnHoldCustomerID := incomingOnHoldCustomerID.(string) // Type assertion
			incomingRequest.OnHoldCustomerID = utils.ToPtr(incomingOnHoldCustomerID)
		}

		if incomingCheckedOutCustomerID, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]; hasCheckedOutCustomerID {
			incomingCheckedOutCustomerID := incomingCheckedOutCustomerID.(string) // Type assertion
			incomingRequest.CheckedOutCustomerID = utils.ToPtr(incomingCheckedOutCustomerID)
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
