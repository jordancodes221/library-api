package handlers

import (
	"bytes"
	"testing"
	"encoding/json"
	"time"
	"example/library_project/utils"
	"example/library_project/models"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"

	"fmt"
)

func TestBooksHandler_UpdateBook(t *testing.T) {
	existingBook1 := &models.Book{
		ISBN: utils.ToPtr("00001"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: nil, 
		TimeUpdated: nil,
	}

	library := map[string]*models.Book{
		"00001": existingBook1,
	}

	h := &BooksHandler{Books: library}

	arbitraryTimeCreated := time.Date(2023, 1, 1, 1, 30, 0, 0, time.UTC)
	fmt.Println(arbitraryTimeCreated)

	tests := []struct{
		description string
		currentBook *models.Book
		incomingBook *models.Book
		expectedStatusCode int
		expectedBook *models.Book
		expectedError *models.ErrorResponse
	}{
		{
			description: "Successfully check out a book",
			currentBook: existingBook1,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00001"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("02"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00001"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("02"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedError: nil,
		},
	}

	r := gin.Default()
	r.PATCH("/books/:isbn", h.UpdateBook)

	for _, currentTestCase := range tests {
		fmt.Println(currentTestCase.description)
		t.Log(currentTestCase.description)

		bookJSON, _ := json.Marshal(*currentTestCase.incomingBook)

		currentISBN := *currentTestCase.currentBook.ISBN

		req, err := http.NewRequest("PATCH", "/books/"+currentISBN, bytes.NewBuffer(bookJSON))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, currentTestCase.expectedStatusCode, w.Code)

		if currentTestCase.expectedBook != nil {
			// Decode response body into Book struct
			 actualBook := new(models.Book)
			 dec := json.NewDecoder(w.Body)
			 if err := dec.Decode(&actualBook); err != nil {
				t.Fatal(err)
			 }

			 // Check if actual book fields are equal to expected
			 // Note we cannot check TimeUpdated as this is set by the handler at run-time
			 assert.Equal(t, currentTestCase.expectedBook.ISBN, actualBook.ISBN)
			 assert.Equal(t, currentTestCase.expectedBook.State, actualBook.State)
			 assert.Equal(t, currentTestCase.expectedBook.OnHoldCustomerID, actualBook.OnHoldCustomerID)
			 assert.Equal(t, currentTestCase.expectedBook.CheckedOutCustomerID, actualBook.CheckedOutCustomerID)
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