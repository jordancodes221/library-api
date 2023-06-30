package handlers

import (
	"testing"
	"encoding/json"
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

func TestBooksHandler_GetAllBooks(t *testing.T) {
	arbitraryTime := time.Date(2023, 2, 1, 1, 30, 0, 0, time.UTC)

	existingBook1 := &models.Book{
		ISBN: utils.ToPtr("00001"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTime), 
		TimeUpdated: nil,
	}

	existingBook2 := &models.Book{
		ISBN: utils.ToPtr("00002"), 
		State: utils.ToPtr("checked-out"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: utils.ToPtr("02"), 
		TimeCreated: utils.ToPtr(arbitraryTime), 
		TimeUpdated: nil,
	}

	existingBook3 := &models.Book{
		ISBN: utils.ToPtr("00003"), 
		State: utils.ToPtr("on-hold"), 
		OnHoldCustomerID: utils.ToPtr("04"), 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTime), 
		TimeUpdated: nil,
	}

	daoFactory := inmemorydao.NewInMemoryDAOFactory()
	bookDAO := daoFactory.BookDAO()
	fixedTimeProvider := &utils.TestingDateTimeProvider{
		ArbitraryTime: arbitraryTime,
	}

	h := NewBooksHandler(bookDAO, fixedTimeProvider)
	h.BookDAOInterface.Create(existingBook1)
	h.BookDAOInterface.Create(existingBook2)
	h.BookDAOInterface.Create(existingBook3)

	tests := []struct{
		description string
		expectedStatusCode int
		expectedBooks *[]models.Book
		expectedError *models.ErrorResponse
	}{
		{
			description: "Successfully get all books",
			expectedStatusCode: 200,
			expectedBooks: &[]models.Book{
				models.Book{
					ISBN: utils.ToPtr("00001"),
					State: utils.ToPtr("available"),
					OnHoldCustomerID: nil,
					CheckedOutCustomerID: nil,
					TimeCreated: utils.ToPtr(arbitraryTime),
					TimeUpdated: nil,
				},
				models.Book{
					ISBN: utils.ToPtr("00002"),
					State: utils.ToPtr("checked-out"),
					OnHoldCustomerID: nil,
					CheckedOutCustomerID: utils.ToPtr("02"),
					TimeCreated: utils.ToPtr(arbitraryTime),
					TimeUpdated: nil,
				},
				models.Book{
					ISBN: utils.ToPtr("00003"),
					State: utils.ToPtr("on-hold"),
					OnHoldCustomerID: utils.ToPtr("04"),
					CheckedOutCustomerID: nil,
					TimeCreated: utils.ToPtr(arbitraryTime),
					TimeUpdated: nil,
				},
			},
			expectedError: nil,
		},
	}

	r := gin.Default()
	r.GET("/books", h.GetAllBooks)

	for _, currentTestCase := range tests {
		fmt.Println(currentTestCase.description)
		t.Log(currentTestCase.description)

		req, err := http.NewRequest("GET", "/books", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, currentTestCase.expectedStatusCode, w.Code)

		if currentTestCase.expectedBooks != nil {
			// Decode response body into Book struct
			actualBooks := new([]models.Book)
			dec := json.NewDecoder(w.Body)
			if err := dec.Decode(&actualBooks); err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, *currentTestCase.expectedBooks, *actualBooks)
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