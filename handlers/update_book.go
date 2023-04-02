package handlers

import (
	"example/library_project/models"
	"example/library_project/utils"
	
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	"encoding/json"
	"fmt"
)

var invalidRequestErr = errors.New("invalid request")
var conflictErr = errors.New("conflict")

// validateLogicForUpdateBook validates requests for the logic unique to updating an existing book
func validateLogicForUpdateBook(incomingBook *models.Book, currentBook *models.Book) (error) {	
	// Ensure ISBN is provided
	if incomingBook.ISBN == nil {
		return fmt.Errorf("Expected 'isbn' to be non-null: %w", invalidRequestErr)
		// return errors.New("Missing ISBN in the incoming request.")
	}
	
	// Ensure state is provided
	if incomingBook.State == nil {
		return fmt.Errorf("Expected 'state' to be non-null: %w", invalidRequestErr)
		// return errors.New("Missing State in the incoming request.")
	}

	// Validate Time Created
	if incomingBook.TimeCreated != nil {
		// Since incomingBook.TimeCreated is not nil, we can de-reference it
		incomingTimeCreated := *incomingBook.TimeCreated

		currentTimeCreated := *currentBook.TimeCreated // should not need to check that currentBook.TimeCreated != nil, because all books have a Time Created and this field cannot be changed by the client
		
		if incomingTimeCreated != currentTimeCreated {
			return fmt.Errorf("'timecreated' cannot be modified: %w", invalidRequestErr)
			// return errors.New("Requested time created does not match existing time created.")
		}
	}

	// Validate Time Updated
	if incomingBook.TimeUpdated != nil {
		incomingTimeUpdated := *incomingBook.TimeUpdated // since incomingBook.TimeUpdated is not nil, we can de-reference it
		// perhaps de-referencing incomingBook.TimeUpdated could be moved to the else-block below
		// However, I am keeping it here so it can be de-refenced on the line after checking it is not nil

		// Now, we check whether the current book has a time updated provided or not
		if currentBook.TimeUpdated == nil {
			return fmt.Errorf("'timeupdated' cannot be modified: %w", invalidRequestErr)
			// return errors.New("Requested time updated does not match existing time updated.")
		} else { // currentBook.TimeUpdated != nil
			currentTimeUpdated := *currentBook.TimeUpdated // in this case, currentBook.TimeUpdated is not nil so we can de-reference it
			if incomingTimeUpdated != currentTimeUpdated {
				return fmt.Errorf("'timeupdated' cannot be modified: %w", invalidRequestErr)
				// return errors.New("Requested time updated does not match existing time updated.")
			}
		}
		
	}

	return nil
}

// validateIDsForCheckedOut ensures the OnHoldCustomerID and CheckedOutCustomerID fields are correctly populated for the checkout and returnBook helper functions
func validateIDsForCheckedOut(incomingBook *models.Book, currentBook *models.Book) (error) {
	if (incomingBook.CheckedOutCustomerID == nil) {
		return fmt.Errorf("Expected 'checkedoutcustomerid' to be non-null: %w", invalidRequestErr)
		// return errors.New("Expected checked-out customer ID.")
	}

	if (incomingBook.OnHoldCustomerID != nil) {
		return fmt.Errorf("Expected 'onholdcustomerid' to be null: %w", invalidRequestErr)
		// return errors.New("Did not expect on-hold customer ID.")
	}

	return nil
}

// validateIDsForOnHold ensures the OnHoldCustomerID and CheckedOutCustomerID fields are correctly populated for the placeHold and releaseHold helper functions
func validateIDsForOnHold(incomingBook *models.Book, currentBook *models.Book) (error) {
	if (incomingBook.OnHoldCustomerID == nil) {
		return fmt.Errorf("Expected 'onholdcustomerid' to be non-null: %w", invalidRequestErr)
		// return errors.New("Expected on-hold customer ID.")
	}

	if (incomingBook.CheckedOutCustomerID != nil) {
		return fmt.Errorf("Expected 'checkedoutcustomerid' to be null: %w", invalidRequestErr)
		// return errors.New("Did not expect checked-out customer ID.")
	}

	return nil
}

// checkout
	// available --> checked-out
	// on-hold --> checked-out
	// checked-out --> checked-out
func checkout(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := validateIDsForCheckedOut(incomingBook, currentBook); err != nil {
		return nil, err
	}

	if (*currentBook.State == "available") {
		*currentBook.State = "checked-out"
		currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
		currentBook.TimeUpdated = utils.ToPtr(time.Now())
	} else if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			*currentBook.State = "checked-out"
			currentBook.OnHoldCustomerID = nil
			currentBook.CheckedOutCustomerID = incomingBook.CheckedOutCustomerID
			currentBook.TimeUpdated = utils.ToPtr(time.Now())
		} else {
			return nil, fmt.Errorf("Checkout failed as another customer has the book on-hold: %w", conflictErr)
			// return nil, errors.New("Cannot complete checkout. Someone else has the book on-hold.")
		}
	} else if (*currentBook.State == "checked-out") {
		if (*currentBook.CheckedOutCustomerID == *incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it checked out is the same one trying to check it out redundantly
			// pass
		} else {
			return nil, fmt.Errorf("Checkout failed as another customer has the book checked-out: %w", conflictErr)
			// return nil, errors.New("Cannot complete checkout. Someone else has the book checked-out.")
		}
	} else {
		// pass
	}

	return currentBook, nil
}

// conflict
	// checked-out --> on-hold
func conflict(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	return nil, fmt.Errorf("Invalid state transition requested: %w", conflictErr)
	// return nil, errors.New("Invalid state transition requested.")
}

// placeHold
	// available --> on-hold
	// on-hold --> on-hold
func placeHold(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := validateIDsForOnHold(incomingBook, currentBook); err != nil {
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
			return nil, fmt.Errorf("Placing hold failed as another customer has the book on-hold: %w", conflictErr)
			// return nil, errors.New("Cannot place hold. Someone else already has the book on-hold.")
		}
	} else {
		// pass 
	}

	return currentBook, nil
}

// releaseHold
	// on-hold --> available
func releaseHold(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := validateIDsForOnHold(incomingBook, currentBook); err != nil {
		return nil, err
	}

	if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.OnHoldCustomerID) {
			*currentBook.State = "available"
			currentBook.OnHoldCustomerID = nil
			currentBook.TimeUpdated = utils.ToPtr(time.Now())
		} else {
			return nil, fmt.Errorf("Releasing hold failed as it is another customer who has the book on-hold: %w", conflictErr)
			// return nil, errors.New("Someone else has this book on hold. You cannot release the hold on a book that do not currently have on-hold.")
		}
	}

	return currentBook, nil
}

// returnBook
	// checked-out --> available
func returnBook(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := validateIDsForCheckedOut(incomingBook, currentBook); err != nil {
		return nil, err
	}

	if (*currentBook.State == "checked-out") {
		if (*currentBook.CheckedOutCustomerID == *incomingBook.CheckedOutCustomerID) {
			*currentBook.State = "available"
			currentBook.CheckedOutCustomerID = nil
			currentBook.TimeUpdated = utils.ToPtr(time.Now())
		} else {
			return nil, fmt.Errorf("Returning the book failed as it is another customer who has the book checked-out: %w", conflictErr)
			// return nil, errors.New("Someone else has this book checked-out. You cannot return a book that you did not check out.")
		}
	}

	return currentBook, nil
}

// noOperation
	// available --> available
	// on-hold --> on-hold (when ID's match)
func noOperation(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	return currentBook, nil
}

var actionTable = map[string]map[string]func(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	"available": {
		"available": noOperation,
		"checked-out": checkout,
		"on-hold": placeHold,
	}, "checked-out": {
			"available": returnBook,
			"checked-out": checkout,
			"on-hold": conflict,
	}, "on-hold": {
			"available": releaseHold,
			"checked-out": checkout,
			"on-hold": placeHold,
	},
}

// UpdateBook allows the client to update the state of an existing book in the library
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

	// Decode JSON to book struct
	incomingBook := new(models.Book) // the "new" keyword allocates memory for models.Book, and returns a pointer to it
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(incomingBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// If fields are not nil, ensure they are within range
	if err := incomingBook.Validate(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Validate logic
	if err := validateLogicForUpdateBook(incomingBook, currentBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Now we will pass the current state and incoming state to the action table
	currentState := currentBook.State // this is a pointer

	incomingState := *incomingBook.State  // due to validateLogicForUpdateBook, we know incomingBook.State is not nil so we can de-reference it

	currentBook, err = actionTable[*currentState][incomingState](currentBook, incomingBook)
	if err != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, currentBook)
}