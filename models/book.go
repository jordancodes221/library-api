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
	if incomingBook.ISBN == nil {
		return errors.New("Missing ISBN in the incoming request.")
	}
	
	// Ensure state is provided
	if incomingBook.State == nil {
		return errors.New("Missing State in the incoming request.")
	}

	// Validate Time Created
	if incomingBook.TimeCreated != nil {
		// Since incomingBook.TimeCreated is not nil, we can de-reference it
		incomingTimeCreated := *incomingBook.TimeCreated

		currentTimeCreated := *currentBook.TimeCreated // should not need to check that currentBook.TimeCreated != nil, because all books have a Time Created and this field cannot be changed by the client
		
		if incomingTimeCreated != currentTimeCreated {
			return errors.New("Requested time created does not match existing time created.")
		}
	}

	// Validate Time Updated
	if incomingBook.TimeUpdated != nil {
		incomingTimeUpdated := *incomingBook.TimeUpdated // since incomingBook.TimeUpdated is not nil, we can de-reference it
		// perhaps de-referencing incomingBook.TimeUpdated could be moved to the else-block below
		// However, I am keeping it here so it can be de-refenced on the line after checking it is not nil

		// Now, we check whether the current book has a time updated provided or not
		if currentBook.TimeUpdated == nil {
			return errors.New("Requested time updated does not match existing time updated.") // should message be more specific?
		} else { // currentBook.TimeUpdated != nil
			currentTimeUpdated := *currentBook.TimeUpdated // in this case, currentBook.TimeUpdated is not nil so we can de-reference it
			if incomingTimeUpdated != currentTimeUpdated {
				return errors.New("Requested time updated does not match existing time updated.") // should message be more specific?
			}
		}
		
	}

	return nil
}

// ValidateIDsForCheckedOut ensures the OnHoldCustomerID and CheckedOutCustomerID fields are correctly populated for the checkout and returnBook helper functions
func (incomingBook *Book) ValidateIDsForCheckedOut(currentBook *Book) (error) {
	if (incomingBook.CheckedOutCustomerID == nil) {
		return errors.New("Expected checked-out customer ID.")
	}

	if (incomingBook.OnHoldCustomerID != nil) {
		return errors.New("Did not expect on-hold customer ID.")
	}

	return nil
}

// ValidateIDsForOnHold ensures the OnHoldCustomerID and CheckedOutCustomerID fields are correctly populated for the placeHold and releaseHold helper functions
func (incomingBook *Book) ValidateIDsForOnHold(currentBook *Book) (error) {
	if (incomingBook.OnHoldCustomerID == nil) {
		return errors.New("Expected on-hold customer ID.")
	}

	if (incomingBook.CheckedOutCustomerID != nil) {
		return errors.New("Did not expect checked-out customer ID.")
	}

	return nil
}