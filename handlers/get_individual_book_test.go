package handlers

import (
	"encoding/json"
	"example/library_project/dao/inmemorydao"
	"example/library_project/models"
	"example/library_project/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"fmt"
	"log"
)

func TestBooksHandler_GetIndividualBook(t *testing.T) {
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

	if err := daoFactory.Open(); err != nil {
		log.Fatal("failed to open database connection: ", err)
	}
	defer daoFactory.Close()

	if err := daoFactory.Clear(); err != nil {
		log.Fatal("failed to clear database: ", err)
	}

	bookDAO := daoFactory.BookDAO()

	bookDAO.Create(existingBook1)

	fixedTimeProvider := &utils.TestingDateTimeProvider{
		ArbitraryTime: arbitraryTime,
	}

	h := NewBooksHandler(bookDAO, fixedTimeProvider)


	tests := []struct{
		description string
		isbn string
		expectedStatusCode int
		expectedBook *models.Book
		expectedError *models.ErrorResponse
	}{
		{
			description: "Successfully get the book with isbn 00001",
			isbn: "00001",
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00001"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: nil,
				TimeCreated: utils.ToPtr(arbitraryTime),
				TimeUpdated: nil,
			},
			expectedError: nil,
		},
		{
			description: "Book not found",
			isbn: "00002",
			expectedStatusCode: 404,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("REQUEST SUCCESSFUL. BOOK NOT FOUND"),
			},
		},
	}

	r := gin.Default()
	r.GET("/books/:isbn", h.GetIndividualBook)

	for _, currentTestCase := range tests {
		fmt.Println(currentTestCase.description)
		t.Log(currentTestCase.description)

		req, err := http.NewRequest("GET", "/books/"+currentTestCase.isbn, nil)
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

			assert.Equal(t, currentTestCase.expectedBook, actualBook)
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