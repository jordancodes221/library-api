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

//
// TODO: The following syntax validation function must take a pointer-receiver of the book struct - This is challenging because we must do it on the JSON
//


// // THIS IS A POINTER-RECEIVER FUNCTION
// // Syntactic Validation
// func (b *Book) ValidateISBNAndStateSyntax() (error) { // the request will not complete if input is not OK, which why it is possible to return an error		
// 	// Assuming ISBN is present, is it valid?
// 	if b.ISBN != nil {
// 		isbn := *b.ISBN
		
// 		// Check it's a string here?
// 			// Interface type assertion (*b.ISBN).(string) will not work... maybe consider a try-catch

// 		if isbn == "" {
// 			return errors.New("ISBN cannot be an empty string.")
// 		}
// 	}

// 	if b.State != nil {
// 		state := *b.State

// 		// Check it's a string here?
// 			// Interface type assertion (*b.ISBN).(string) will not work... maybe consider a try-catch

// 		if ((state != "available") && (state != "on-hold") && (state != "checked-out")) {
// 			return errors.New("Invalid state provided. State must be equal to one of: \"available\", \"on-hold\", or \"checked-out\".")
// 			// Tested in "Invalid State" test of UpdateBook in Postman
// 		}

// 	}
	
// 	return nil
// }

// NOT A POINTER-RECEIVER FUNCTION
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