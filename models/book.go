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