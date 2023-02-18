package handlers

import ( // ValidateTimeSemanticsForCreateBook
	"example/library_project/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"encoding/json"
	"errors"
)

// Validate Time Semantics
func ValidateTimeSemanticsForCreateBook(incomingBookAsMap map[string]interface{}) (error) {
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

// POST
func CreateBook(c *gin.Context) {
	var newBook *models.Book = &models.Book{
		ISBN: nil, 
		State: nil, 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: nil, 
		TimeUpdated: nil,}

	// Unmarshal
	incomingBookAsMap := map[string]interface{}{}
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Ensure that incoming JSON includes ISBN
	if _, hasISBN := incomingBookAsMap["isbn"]; !hasISBN {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Missing ISBN in the incoming request."})
		return
	}

	// The reason for calling validate at this point is that it is:
		// (1) After checking that ISBN is present, and
		// (2) Before checking if ISBN is in-use (we want to ensure it's valid before checking if it's in-use)

	// Validate ISBN and State Syntax
	if err := models.ValidateISBNAndStateSyntax(incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Make sure ISBN is not already in-use
		// At this point, we know that ISBN (1) is present, and (2) is valid
	incomingISBN := incomingBookAsMap["isbn"].(string)
	if _, ok := mapOfBooks[incomingISBN]; ok {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": "Book already exists."})
		return
	}

	// Validate Time Semantics
	if err := ValidateTimeSemanticsForCreateBook(incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Update newBook ISBN field
	newBook.ISBN = ToPtr(incomingISBN)

	// Update newBook State field (if state is present)
	if incomingState, hasState := incomingBookAsMap["state"]; hasState {
		incomingState := incomingState.(string) // Type Assertion
		newBook.State = ToPtr(incomingState)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Missing State in the incoming request."})
		return
	}

	// Ensure correct customer ID fields are provided given the state
	if err := ValidateIDSemanticsForCreateBook(incomingBookAsMap); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Update newBook OnHoldCustomerID
	if incomingOnHoldCustomerID, hasOnHoldCustomerID := incomingBookAsMap["onholdcustomerid"]; hasOnHoldCustomerID {
		incomingOnHoldCustomerID := incomingOnHoldCustomerID.(string) // Type assertion
		newBook.OnHoldCustomerID = ToPtr(incomingOnHoldCustomerID)
	}

	// Update newBook CheckedOutCustomerID
	if incomingCheckedOutCustomerID, hasCheckedOutCustomerID := incomingBookAsMap["checkedoutcustomerid"]; hasCheckedOutCustomerID {
		incomingCheckedOutCustomerID := incomingCheckedOutCustomerID.(string) // Type assertion
		newBook.CheckedOutCustomerID = ToPtr(incomingCheckedOutCustomerID)
	}

	// Update newBook times
	newBook.TimeCreated = ToPtr(time.Now())
	newBook.TimeUpdated = nil

	// Add newBook to mapOfBooks
	mapOfBooks[*newBook.ISBN] = newBook

	c.IndentedJSON(http.StatusCreated, newBook) // 201 status code if successful
}