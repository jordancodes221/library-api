package handlers

import (
	"example/library_project/models"
	// "example/library_project/utils"

	"net/http"
	"github.com/gin-gonic/gin"
	// "time"
	"encoding/json"
	"errors"
)

// validateLogicForCreateBook validates requests for the logic specific to creating a new book
func validateLogicForCreateBook(incomingBook *models.Book) (error) {
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

// CreateBook allows the client to add a new book to the library
func (h *BooksHandler) CreateBook(c *gin.Context) {
	// Decode JSON to book struct
	newBook := new(models.Book) // the "new" keyword allocates memory for models.Book, and returns a pointer to it
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// If fields are not nil, ensure they are within range
	if err := newBook.Validate(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Logic validation
	if err := validateLogicForCreateBook(newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Make sure ISBN is not already in-use
	bookWithISBNInUse, err := h.BookDAOInterface.Read(*newBook.ISBN)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if bookWithISBNInUse != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": "Book already exists."})
		return
	}

	// Update TimeCreated to now
	newBook.TimeCreated = h.DateTimeInterface.GetCurrentTime()

	// Add the new book to our library
	if err := h.BookDAOInterface.Create(newBook); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newBook) // 201 status code if successful
}