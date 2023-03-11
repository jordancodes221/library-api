package handlers

import ( // h.books, h.bookByISBN
	"example/library_project/models"
	"example/library_project/utils"
	
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"time"
	"encoding/json"
)

// No Match Error
var noMatchError error = errors.New("ID's do not match.")

// checkout
	// available --> checked-out
	// on-hold --> checked-out
	// checked-out --> checked-out
func checkout(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
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
			return nil, noMatchError
		}
	} else if (*currentBook.State == "checked-out") {
		if (*currentBook.CheckedOutCustomerID == *incomingBook.CheckedOutCustomerID) { // ensure the customer who currently has it checked out is the same one trying to check it out redundantly
			// pass
		} else {
			return nil, noMatchError
		}
	} else {
		// pass
	}

	return currentBook, nil
}

// conflict
	// checked-out --> on-hold
func conflict(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	return nil, errors.New("Invalid state transfer requested.")
}

// placeHold
	// available --> on-hold
	// on-hold --> on-hold
func placeHold(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if (*currentBook.State == "available") {
		*currentBook.State = "on-hold"
		currentBook.OnHoldCustomerID = incomingBook.OnHoldCustomerID
		currentBook.TimeUpdated = utils.ToPtr(time.Now())
	} else if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.OnHoldCustomerID) { // ensure the customer who currently has it on-hold is the same one trying to check it out
			// pass
		} else {
			return nil, noMatchError
		}
	} else {
		// pass 
	}

	return currentBook, nil
}

// releaseHold
	// on-hold --> available
func releaseHold(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if (*currentBook.State == "on-hold") {
		if (*currentBook.OnHoldCustomerID == *incomingBook.OnHoldCustomerID) {
			*currentBook.State = "available"
			currentBook.OnHoldCustomerID = nil
			currentBook.TimeUpdated = utils.ToPtr(time.Now())
		} else {
			return nil, noMatchError
		}
	}

	return currentBook, nil
}

// returnBook
	// checked-out --> available
func returnBook(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if (*currentBook.State == "checked-out") {
		if (*currentBook.CheckedOutCustomerID == *incomingBook.CheckedOutCustomerID) {
			*currentBook.State = "available"
			currentBook.CheckedOutCustomerID = nil
			currentBook.TimeUpdated = utils.ToPtr(time.Now())
		} else {
			return nil, noMatchError
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

// First key is current state, 2nd key is incoming state
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

	// General validation for logic
	if err := incomingBook.GeneralValidationForUpdateBook(currentBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Now we will pass the current state and incoming state to the action table
	currentState := currentBook.State // this is a pointer

	ptrIncomingState := incomingBook.State 
	incomingState := *ptrIncomingState // due to GeneralValidationForUpdateBook, we know ptrIncomingState is not nil

	currentBook, err = actionTable[*currentState][incomingState](currentBook, incomingBook) // now the action table needs to be re-written to 2nd parameter is struct instead of map
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, currentBook)
}