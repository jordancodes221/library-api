package models

import (
	"time"
	"errors"
)

// Book represents an individual book in the library
type Book struct{
	// ISBN is a unique identifier for the book
	ISBN 			*string 	`json:"isbn"`

	// State is the current state of the book. It can be "available", "on-hold", or "checked-out"
	State 			*string 	`json:"state"`

	// OnHoldCustomerID identifies the customer who has the book on-hold. This field must also be provided in any request to place or release a hold on a book
	OnHoldCustomerID 	*string 	`json:"onholdcustomerid"`

	// CheckedOutCustomerID identifies the customer who has the book checked-out. This field must also be provided in any request to checkout or return a book
	CheckedOutCustomerID 	*string 	`json:"checkedoutcustomerid"`

	// TimeCreated is the time the book was created. It is immutable by the client
	TimeCreated 		*time.Time 	`json:"timecreated"`

	// TimeUpdated is the time the book was last updated. It is immutable by the client
	TimeUpdated  		*time.Time	`json:"timeupdated"`
}

// Validate ensures that all fields provided in the request are within range for both creating a new book and updating an existing book
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

// ValidateLogicForCreateBook validates requests for the logic specific to creating a new book
func (incomingBook *Book) ValidateLogicForCreateBook() (error) {
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

// ValidateLogicForUpdateBook validates requests for the logic unique to updating an existing book
func (incomingBook *Book) ValidateLogicForUpdateBook(currentBook *Book) (error) {	
	// Ensure ISBN is provided
	ptrIncomingISBN := incomingBook.ISBN
	if ptrIncomingISBN == nil {
		return errors.New("Missing ISBN in the incoming request.")
	}
	
	// Ensure state is provided
	ptrIncomingState := incomingBook.State
	if ptrIncomingState == nil {
		return errors.New("Missing State in the incoming request.")
	}

	// Ensure on-hold and checked-out IDs are provided correctly, given the state
	currentState := *currentBook.State
	incomingState := *ptrIncomingState // at this point we know it is not nil, so we can de-reference it

	ptrCheckedOutCustomerID := incomingBook.CheckedOutCustomerID
	ptrOnHoldCustomerID := incomingBook.OnHoldCustomerID

	// This corresponds to the checkout and returnBook helper functions in the action table
	if ((currentState == "available" && incomingState == "checked-out") || (currentState == "checked-out" && incomingState == "available")){	
		if (ptrCheckedOutCustomerID == nil) {
			return errors.New("Expected checked-out customer ID.")
		}
	
		if (ptrOnHoldCustomerID != nil) {
			return errors.New("Did not expect on-hold customer ID.")
		}
	}

	// This corresponds to the placeHold and releaseHold helper functions in the action table
	if ((currentState == "available" && incomingState == "on-hold") || (currentState == "on-hold" && incomingState == "available")){
		if (ptrOnHoldCustomerID == nil) {
			return errors.New("Expected on-hold customer ID.")
		}
	
		if (ptrCheckedOutCustomerID != nil) {
			return errors.New("Did not expect checked-out customer ID.")
		}
	}

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