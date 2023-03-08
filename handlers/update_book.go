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
func validateTimeSemanticsForUpdateBook(currentBook *models.Book, incomingBookAsMap map[string]interface{}) (error) {
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


// Semantic Validation for checkout and returnBook
func validateIDSemanticsForCheckedOutUpdate(incomingRequest *models.Book) (error) {
	// incomingRequest is of the form &{isbn, state, checkedoutcustomerid, onholdcustomerid, timecreated, timeupdated}
	// For this particular case, it should be populated as such: &{isbn, state, checkedoutcustomerid, nil, nil, nil}
	
	// fmt.Println("CALLING validateIDSemanticsForCheckedOutUpdate...")
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

// Semantic Validation for placeHold and releaseHold
func validateIDSemanticsForOnHoldUpdate(incomingRequest *models.Book) (error) {
	// incomingRequest is of the form &{isbn, state, checkedoutcustomerid, onholdcustomerid, timecreated, timeupdated}
	// For this particular case, it should be populated as such: &{isbn, state, nil, onholdcustomerid, nil, nil}
	checkedOutCustomerID := incomingRequest.CheckedOutCustomerID
	onHoldCustomerID := incomingRequest.OnHoldCustomerID

	// fmt.Println("CALLING validateIDSemanticsForOnHoldUpdate...")
	if (onHoldCustomerID == nil) {
		return errors.New("Expected on-hold customer ID.")
	}

	if (checkedOutCustomerID != nil) {
		return errors.New("Did not expect checked-out customer ID.")
	}

	return nil
}

// No Match Error
var noMatchError error = errors.New("ID's do not match.")

// checkout
	// available --> checked-out
	// on-hold --> checked-out
	// checked-out --> checked-out
func checkout(currentBook *models.Book, incomingBook *models.Book) (*models.Book, error) {
	if err := validateIDSemanticsForCheckedOutUpdate(incomingBook); err != nil {
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
	if err := validateIDSemanticsForOnHoldUpdate(incomingBook); err != nil {
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
	if err := validateIDSemanticsForOnHoldUpdate(incomingBook); err != nil {
		return nil, err
	}

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
	fmt.Println("CALLING RETURNBOOK...")

	if err := validateIDSemanticsForCheckedOutUpdate(incomingBook); err != nil {
		return nil, err
	}

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

///////////////////////////////////
/////// MAJOR CHANGES BEGIN ///////
///////////////////////////////////

	// Decode JSON to book struct
			// Notice this block is the same as the previous block with decoding to a map, except the first line
			// Previously, the first line allocated memory for a map[string]interface{}{}, but here we allocae memory for a models.Book struct
			// Also, in the 3rd line we pass book (rather than a pointer to the map)... note that the book variable is a pointer (see commment on that line)
	incomingBookAsStruct := new(models.Book) // the "new" keyword allocates memory for models.Book, and returns a pointer to it
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(incomingBookAsStruct); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// If fields are not nil, ensure they are within range
	if err := incomingBookAsStruct.Validate(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// General validation for logic (additional validation needed depending on specific action table helper function called later)
	if err := incomingBookAsStruct.GeneralValidationForUpdateBook(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	//// NOW WE GET OUR CURRENT AND INCOMING STATES, SO WE CAN PASS THEM INTO THE ACTION TABLE

	// this remains unchanged, as it is the current book (not the incoming book)
	currentState := currentBook.State // this is a pointer

	ptrIncomingState := incomingBookAsStruct.State 
	incomingState := *ptrIncomingState // due to GeneralValidationForUpdateBook, we know ptrIncomingState is not nil

	currentBook, err = actionTable[*currentState][incomingState](currentBook, incomingBookAsStruct) // now the action table needs to be re-written to 2nd parameter is struct instead of map
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, currentBook)
}