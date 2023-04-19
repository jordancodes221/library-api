package handlers

import (
	"bytes"
	"testing"
	"encoding/json"
	// "time"
	"example/library_project/utils"
	"example/library_project/models"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"

	"fmt"
)

func TestBooksHandler_CreateBook(t *testing.T) {
	existingBook := &models.Book{
		ISBN: utils.ToPtr("11111"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: nil, 
		TimeUpdated: nil,
	}

	library := map[string]*models.Book{
		"11111": existingBook,
	}

	h := &BooksHandler{Books: library}
	
	tests := []struct{
		description string
		book *models.Book
		expectedStatusCode int
		expectedBook *models.Book
		expectedError *models.ErrorResponse
	}{
		{
			description: "Valid book", 
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 201,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			},
			expectedError: nil,
		}, 
		{
			description: "ISBN is the empty string", 
			book: &models.Book{
				ISBN: utils.ToPtr(""), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("ISBN cannot be the empty string."),
			},
		},
	}
	
	for _, currentTestCase := range tests {
		fmt.Println(currentTestCase.description)
		t.Log(currentTestCase.description)

		bookJSON, _ := json.Marshal(*currentTestCase.book)
		fmt.Println(*currentTestCase.book)
		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJSON))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		h.CreateBook(c)

		assert.Equal(t, currentTestCase.expectedStatusCode, w.Code)

		if currentTestCase.expectedBook != nil {
			// Decode response body into Book struct
			 actualBook := new(models.Book)
			 dec := json.NewDecoder(w.Body)
			 if err := dec.Decode(&actualBook); err != nil {
				t.Fatal(err)
			 }

			 // Check if actual book fields are equal to expected
			 // Note we cannot check TimeCreated as this is set by the handler at run-time
			 assert.Equal(t, currentTestCase.expectedBook.ISBN, actualBook.ISBN)
			 assert.Equal(t, currentTestCase.expectedBook.State, actualBook.State)
			 assert.Equal(t, currentTestCase.expectedBook.OnHoldCustomerID, actualBook.OnHoldCustomerID)
			 assert.Equal(t, currentTestCase.expectedBook.CheckedOutCustomerID, actualBook.CheckedOutCustomerID)
			 assert.Equal(t, currentTestCase.expectedBook.TimeUpdated, actualBook.TimeUpdated)
		}

		if currentTestCase.expectedError != nil {
			// Decode response body into ErrorResponse struct
			actualError := new(models.ErrorResponse)
			dec := json.NewDecoder(w.Body)
			if err := dec.Decode(&actualError); err != nil {
				t.Fatal(err)
			}

			// Check if actual error is equal to expected
			assert.Equal(t, currentTestCase.expectedError, actualError)
		}
	}
}
