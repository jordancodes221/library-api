package models

import (
	"time"
	"errors"
)

type Book struct{
	ISBN 			*string 	`json:"isbn"`
	State 			*string 	`json:"state"`

	OnHoldCustomerID 	*string 	`json:"onholdcustomerid"`
	CheckedOutCustomerID 	*string 	`json:"checkedoutcustomerid"`

	TimeCreated 		*time.Time 	`json:"timecreated"`
	TimeUpdated  		*time.Time	`json:"timeupdated"`
}

// Only if a field is present, we validate it to make sure if it is within range
func (incomingBook *Book) Validate() (error) {

	// ISBN
	if incomingBook.ISBN != nil {
		if *incomingBook.ISBN == "" { 	// Remark: In the first if-statement, we check the pointer to the ISBN field. In the 2nd if-statement, we check its value.
			return errors.New("ISBN cannot be the empty string.")
		}
	}

	// State - Tested in "Invalid State" test of UpdateBook in Postman
	if incomingBook.State != nil {
		if ((*incomingBook.State != "available") && (*incomingBook.State != "on-hold") && (*incomingBook.State != "checked-out")) {
			return errors.New("Invalid state provided. State must be equal to one of: \"available\", \"on-hold\", or \"checked-out\".")
		}
	}

	// OnHoldCustomerID
	if incomingBook.OnHoldCustomerID != nil {
		if *incomingBook.OnHoldCustomerID == "" {
			return errors.New("On-hold customer ID cannot be the empty string.")
		}
	}

	// CheckedOutCustomerID
	if incomingBook.CheckedOutCustomerID != nil {
		if *incomingBook.CheckedOutCustomerID == "" {
			return errors.New("Checked-out customer ID cannot be the empty string.")
		}
	}

	return nil
}

//////////////////////////////// 
////////// CREATE BOOK /////////
////////////////////////////////

func (incomingBook *Book) FurtherValidationForCreateBook() (error) {
	// Ensure ISBN is provided
	if incomingBook.ISBN == nil {
		return errors.New("Missing ISBN in the incoming request.")
	}
	
	// Ensure state is provided
	if incomingBook.State == nil {
		return errors.New("Missing State in the incoming request.")
	}

	// CreateBook calls Validate(), which ensures *incomingBook.State (if provided) is equal to one "available", "on-hold", or "checked-out"
	// Since we have already addressed the case where State is not provided, we know at this point State is equal to one of those 3 values.

	// State is Available
	if (*incomingBook.State == "available") {
		if incomingBook.OnHoldCustomerID != nil {
			return errors.New("Cannot have an on-hold customer ID when state is available.")
		}

		if incomingBook.CheckedOutCustomerID != nil {
			return errors.New("Cannot have checked-out customer ID when state is available.")
		}
	}

	// State is On-Hold
	if (*incomingBook.State == "on-hold") {
		if incomingBook.CheckedOutCustomerID != nil {
			return errors.New("Cannot have checked-out customer ID when state is on-hold.")
		}

		if incomingBook.OnHoldCustomerID == nil {
			return errors.New("State provided is on-hold, but no on-hold customer ID is provided.")
		}
	}

	// State is Checked-Out
	if (*incomingBook.State == "checked-out") {
		if incomingBook.OnHoldCustomerID != nil {
			return errors.New("Cannot have on-hold customer ID when state is checked-out.")
		}

		if incomingBook.CheckedOutCustomerID == nil {
			return errors.New("State provided is checked-out, but no checked-out customer ID is provided.")
		}
	}

	// Ensure TimeCreated is not provided by the client
	if incomingBook.TimeCreated != nil {
		return errors.New("Client cannot provide time created when creating a new book.")
	}

	// Ensure TimeUpdated is not provided by the client
	if incomingBook.TimeUpdated != nil {
		return errors.New("Client cannot provide time updated when creating a new book.")
	}

	return nil
}

//////////////////////////////// 
////////// UPDATE BOOK /////////
//////////////////////////////// 

func (incomingBook *Book) GeneralValidationForUpdateBook() (error) {
	// Ensure ISBN is provided
	if incomingBook.ISBN == nil {
		return errors.New("Missing ISBN in the incoming request.")
	}
	
	// Ensure state is provided
	if incomingBook.State == nil {
		return errors.New("Missing State in the incoming request.")
	}

	return nil
}

// Validate Time Semantics
func (incomingBook *Book) ValidateTimeSemanticsForUpdateBook(currentBook *Book) (error) {
	// Validate Time Created
	ptrIncomingTimeCreated := incomingBook.TimeCreated
	if ptrIncomingTimeCreated != nil {
		// Since ptrIncomingTimeCreated is not nil, we can de-reference it
		incomingTimeCreated := *ptrIncomingTimeCreated

		ptrCurrentTimeCreated := currentBook.TimeCreated
		currentTimeCreated := *ptrCurrentTimeCreated // should not need to check that ptrCurrentTimeCreated != nil, because all books have a Time Created and this field cannot be changed by the client
		
		if incomingTimeCreated != currentTimeCreated {
			return errors.New("Requested time created does not match existing time created.")
		}
	}

	// Validate Time Updated
	ptrIncomingTimeUpdated := incomingBook.TimeUpdated
	if ptrIncomingTimeUpdated != nil {
		incomingTimeUpdated := *ptrIncomingTimeUpdated // since ptrIncomingTimeUpdated is not nil, we can de-reference it
		// perhaps de-referencing ptrIncomingTimeUpdated could be moved to the else-block below
		// However, I am keeping it here so it can be de-refenced on the line after checking it is not nil

		// Now, we check whether the current book has a time updated provided or not
		ptrCurrentTimeUpdated := currentBook.TimeUpdated // keep in mind this could be nil
		if ptrCurrentTimeUpdated == nil {
			return errors.New("Requested time updated does not match existing time updated.") // should message be more specific?
		} else { // ptrCurrentTimeUpdated != nil
			currentTimeUpdated := *ptrCurrentTimeUpdated // in this case, ptrCurrentTimeUpdated is not nil so we can de-reference it
			if incomingTimeUpdated != currentTimeUpdated {
				return errors.New("Requested time updated does not match existing time updated.") // should message be more specific?
			}
		}
		
	}

	return nil
}

// Semantic Validation for checkout and returnBook
func (incomingBook *Book) ValidateIDSemanticsForCheckedOutUpdate() (error) {
	checkedOutCustomerID := incomingBook.CheckedOutCustomerID
	onHoldCustomerID := incomingBook.OnHoldCustomerID

	if (checkedOutCustomerID == nil) {
		return errors.New("Expected checked-out customer ID.")
	}

	if (onHoldCustomerID != nil) {
		return errors.New("Did not expect on-hold customer ID.")
	}

	return nil
}

// Semantic Validation for placeHold and releaseHold
func (incomingBook *Book) ValidateIDSemanticsForOnHoldUpdate() (error) {
	checkedOutCustomerID := incomingBook.CheckedOutCustomerID
	onHoldCustomerID := incomingBook.OnHoldCustomerID

	if (onHoldCustomerID == nil) {
		return errors.New("Expected on-hold customer ID.")
	}

	if (checkedOutCustomerID != nil) {
		return errors.New("Did not expect checked-out customer ID.")
	}

	return nil
}