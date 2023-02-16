package validators

import ( // models.
	"errors"
	"example/library_project/models"
	// "time"
	// "fmt"
	// "reflect"
	// "strconv"
)

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

/////////////////////////////////////////////////////////////////////////////
//////////////////////////// Semantic Validation ////////////////////////////
/////////////////////////////////////////////////////////////////////////////

func ValidateTimeSemantics(incomingBookAsMap map[string]interface{}) (error) {
	// fmt.Println("CALLING VALIDATE TIME SEMANTICS...")

	_, hasTimeCreated := incomingBookAsMap["timecreated"]
	_, hasTimeUpdated := incomingBookAsMap["timeupdated"]

	if (hasTimeCreated && !hasTimeUpdated) {
		return errors.New("Client cannot provide time created.")
	}

	if (!hasTimeCreated && hasTimeUpdated) {
		return errors.New("Client cannot provide time updated.")
	}

	if (hasTimeCreated && hasTimeUpdated) {
		return errors.New("Client cannot provide time created or time updated.")
	}

	return nil
}

func ValidateIDSemanticsForCreateBook(incomingBookAsMap map[string]interface{}) (error) {
	// fmt.Println("CALLING ValidateIDSemanticsForCreateBook...")

	// This function will only be called once state is established to be both present and valid
	state := incomingBookAsMap["state"]
	
	// Retrieve the customer ID's if they are present
	_, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]
	_, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]

	// State is available -- THIS IS SEMANTIC CHECKING
	if (state == "available") {
		if ((hasOnHoldCustomerID) && (!hasCheckedOutCustomerID)) {
			return errors.New("Cannot have an on-hold customer ID when state is available.")
		}

		if (!(hasOnHoldCustomerID) && (hasCheckedOutCustomerID)) {
			return errors.New("Cannot have checked-out customer ID when state is available.")
		}
		
		if (hasOnHoldCustomerID || hasCheckedOutCustomerID) {
			return errors.New("Cannot have on-hold customer ID or checked-out customer ID when state is available.")
		}
	}

	// State is on-hold -- THIS IS SEMANTIC CHECKING
	if (state == "on-hold") {
		if hasCheckedOutCustomerID {
			return errors.New("Cannot have checked-out customer ID when state is on-hold.")
		}

		if hasOnHoldCustomerID {
			// We know ohid is provided. Ensure it is a string
			ohid, ohidIsString := incomingBookAsMap["onholdcustomerid"].(string)
			if !ohidIsString {
				return errors.New("On-hold customer ID provided is not of type string.")
			}

			if (ohid == "") {
				return errors.New("On-hold customer ID is the empty string.")
			}
		} else { // !hasOnHoldCustomerID
			return errors.New("State provided is on-hold, but no on-hold customer ID is provided.")
		}
	}

	// State is checked-out -- THIS IS SEMANTIC CHECKING
	if (state == "checked-out") {
		if hasOnHoldCustomerID {
			return errors.New("Cannot have on-hold customer ID when state is checked-out.")
		}

		if hasCheckedOutCustomerID {
			// We know ohid is provided. Ensure it is a string
			coid, coidIsString := incomingBookAsMap["checkedoutcustomerid"].(string)
			if !coidIsString {
				return errors.New("Checked-out customer ID provided is not of type string.")
			}

			if (coid == "") {
				return errors.New("Checked-out customer ID is the empty string.")
			}
		} else { // !hasCheckedOutCustomerID
			return errors.New("State provided is checked-out, but no checked-out customer ID is provided.")
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

	if (checkedOutCustomerID == nil && onHoldCustomerID == nil) {
		return errors.New("Expected checked-out customer ID.")
	}

	if (checkedOutCustomerID != nil && onHoldCustomerID != nil) {
		return errors.New("Did not expect on-hold customer ID.")
	}

	if (checkedOutCustomerID == nil && onHoldCustomerID != nil) {
		return errors.New("Expected checked-out customer ID, and did not expect on-hold customer ID.")
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
	if (onHoldCustomerID == nil && checkedOutCustomerID == nil) {
		return errors.New("Expected on-hold customer ID.")
	}

	if (onHoldCustomerID != nil && checkedOutCustomerID != nil) {
		return errors.New("Did not expect checked-out customer ID.")
	}

	if (onHoldCustomerID == nil && checkedOutCustomerID != nil) {
		return errors.New("Expected on-hold customer ID, and did not expect checked-out customer ID.")
	}

	return nil
}