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

// NOT A POINTER-RECEIVER FUNCTION - still used in update_book.go
// Syntactic Validation
func ValidateISBNAndStateSyntax(incomingBookAsMap map[string]interface{}) (error) { // the request will not complete if input is not OK, which why it is possible to return an error		
	// Assuming ISBN is present, is it valid?
	if isbn, hasISBN := incomingBookAsMap["isbn"]; hasISBN {
		_, isbnIsString := isbn.(string)
		if !isbnIsString {
			return errors.New("ISBN provided is not of type string.") // Tested in "Incorrect ISBN Type" test of CreateBook in Postman
		}
	}

	// Assuming State is present, is it valid?
	if state, hasState := incomingBookAsMap["state"]; hasState {
		state, stateIsString := state.(string)
		if !stateIsString {
			return errors.New("State provided is not of type string.") // Tested in "Incorrect State Type" test of UpdateBook in Postman
		}

		if ((state != "available") && (state != "on-hold") && (state != "checked-out")) {
			return errors.New("Invalid state provided. State must be equal to one of: \"available\", \"on-hold\", or \"checked-out\".")
			// Tested in "Invalid State" test of UpdateBook in Postman
		}
	}

	return nil
}

// NEW VERSION OF POINTER-RECEIVER FUNCTION ADDED FOR LATER USE
	// The following function takes a pointer-receiver of the book struct - Hence, this validation must be performed on the JSON (not on the book as a map)
	// Only if a field is present, we validate it to make sure if it is within range
func (b *Book) Validate() (error) {

	// ISBN
	if b.ISBN != nil {
		if *b.ISBN == "" { 	// Remark: In the first if-statement, we check the pointer to the ISBN field. In the 2nd if-statement, we check its value.
			return errors.New("ISBN cannot be the empty string.")
		}
	}

	// State - Tested in "Invalid State" test of UpdateBook in Postman
	if b.State != nil {
		if ((*b.State != "available") && (*b.State != "on-hold") && (*b.State != "checked-out")) {
			return errors.New("Invalid state provided. State must be equal to one of: \"available\", \"on-hold\", or \"checked-out\".")
		}
	}

	// OnHoldCustomerID
	if b.OnHoldCustomerID != nil {
		if *b.OnHoldCustomerID == "" {
			return errors.New("On-hold customer ID cannot be the empty string.")
		}
	}

	// CheckedOutCustomerID
	if b.CheckedOutCustomerID != nil {
		if *b.CheckedOutCustomerID == "" {
			return errors.New("Checked-out customer ID cannot be the empty string.")
		}
	}

	return nil
}

func (b *Book) FurtherValidationForCreateBook() (error) {
	// Ensure ISBN is provided
	if b.ISBN == nil {
		return errors.New("Missing ISBN in the incoming request.")
	}
	
	// Ensure state is provided
	if b.State == nil {
		return errors.New("Missing State in the incoming request.")
	}

	// CreateBook calls Validate(), which ensures *b.State (if provided) is equal to one "available", "on-hold", or "checked-out"
	// Since we have already addressed the case where State is not provided, we know at this point State is equal to one of those 3 values.

	// State is Available
	if (*b.State == "available") {
		if b.OnHoldCustomerID != nil {
			return errors.New("Cannot have an on-hold customer ID when state is available.")
		}

		if b.CheckedOutCustomerID != nil {
			return errors.New("Cannot have checked-out customer ID when state is available.")
		}
	}

	// State is On-Hold
	if (*b.State == "on-hold") {
		if b.CheckedOutCustomerID != nil {
			return errors.New("Cannot have checked-out customer ID when state is on-hold.")
		}

		if b.OnHoldCustomerID == nil {
			return errors.New("State provided is on-hold, but no on-hold customer ID is provided.")
		}
	}

	// State is Checked-Out
	if (*b.State == "checked-out") {
		if b.OnHoldCustomerID != nil {
			return errors.New("Cannot have on-hold customer ID when state is checked-out.")
		}

		if b.CheckedOutCustomerID == nil {
			return errors.New("State provided is checked-out, but no checked-out customer ID is provided.")
		}
	}

	// Ensure TimeCreated is not provided by the client
	if b.TimeCreated != nil {
		return errors.New("Client cannot provide time created when creating a new book.")
	}

	// Ensure TimeUpdated is not provided by the client
	if b.TimeUpdated != nil {
		return errors.New("Client cannot provide time updated when creating a new book.")
	}

	return nil
}

func (b *Book) GeneralValidationForUpdateBook() (error) {
	// Ensure ISBN is provided
	if b.ISBN == nil {
		return errors.New("Missing ISBN in the incoming request.")
	}
	
	// Ensure state is provided
	if b.State == nil {
		return errors.New("Missing State in the incoming request.")
	}

	return nil
}