package handlers

import (
	"testing"
	// "encoding/json"
	"time"
	"example/library_project/utils"
	"example/library_project/models"
	"example/library_project/dao/inmemorydao"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"

	"fmt"
)

func TestBooksHandler_DeleteBook(t *testing.T) {
	arbitraryTime := time.Date(2023, 2, 1, 1, 30, 0, 0, time.UTC)

	existingBook1 := &models.Book{
		ISBN: utils.ToPtr("00001"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTime), 
		TimeUpdated: nil,
	}

	daoFactory := inmemorydao.NewInMemoryDAOFactory()
	fixedTimeProvider := &utils.TestingDateTimeProvider{
		ArbitraryTime: arbitraryTime,
	}

	h := NewBooksHandler(daoFactory, fixedTimeProvider)
	h.BookDAOInterface.Create(existingBook1)

	tests := []struct{
		description string
		isbn string
		expectedStatusCode int
		expectedError *models.ErrorResponse
	}{
		{
			description: "Successfully delete a book",
			isbn: "00001",
			expectedStatusCode: 204,
			expectedError: nil,
		},
		{
			description: "Book not found",
			isbn: "00002",
			expectedStatusCode: 204,
			expectedError: nil,
		},
	}

	r := gin.Default()
	r.DELETE("/books/:isbn", h.DeleteBook)

	for _, currentTestCase := range tests {
		fmt.Println(currentTestCase.description)
		t.Log(currentTestCase.description)

		req, err := http.NewRequest("DELETE", "/books/"+currentTestCase.isbn, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, currentTestCase.expectedStatusCode, w.Code)

		if currentTestCase.expectedError == nil {
			assert.Empty(t, w.Body)

			// deletedBook, err := h.BookDAOInterface.Read(currentTestCase.isbn)
			// assert.Nil(t, deletedBook)
		}

		// if currentTestCase.expectedError != nil {
		// 	// Decode response body into ErrorResponse struct
		// 	actualError := new(models.ErrorResponse)
		// 	dec := json.NewDecoder(w.Body)
		// 	if err := dec.Decode(&actualError); err != nil {
		// 		t.Fatal(err)
		// 	}

		// 	// Check if actual error is equal to expected
		// 	assert.Equal(t, currentTestCase.expectedError, actualError)
		// }
	}
}