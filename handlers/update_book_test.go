package handlers

import (
	"bytes"
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

func TestBooksHandler_UpdateBook(t *testing.T) {
	arbitraryTimeCreated := time.Date(2023, 2, 1, 1, 30, 0, 0, time.UTC)
	arbitraryTimeUpdated := time.Date(2023, 2, 2, 1, 30, 0, 0, time.UTC)
	
	existingBook1 := &models.Book{
		ISBN: utils.ToPtr("00001"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook2 := &models.Book{
		ISBN: utils.ToPtr("00002"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook3 := &models.Book{
		ISBN: utils.ToPtr("00003"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}
	
	existingBook4 := &models.Book{
		ISBN: utils.ToPtr("00004"), 
		State: utils.ToPtr("checked-out"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: utils.ToPtr("02"), 
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook5 := &models.Book{
		ISBN: utils.ToPtr("00005"), 
		State: utils.ToPtr("checked-out"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: utils.ToPtr("02"), 
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook6 := &models.Book{
		ISBN: utils.ToPtr("00006"), 
		State: utils.ToPtr("on-hold"), 
		OnHoldCustomerID: utils.ToPtr("04"), 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook7 := &models.Book{
		ISBN: utils.ToPtr("00007"), 
		State: utils.ToPtr("on-hold"), 
		OnHoldCustomerID: utils.ToPtr("06"), 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook8 := &models.Book{
		ISBN: utils.ToPtr("00008"), 
		State: utils.ToPtr("on-hold"), 
		OnHoldCustomerID: utils.ToPtr("04"), 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	// exitingBook9 will be used for the invalid checked-out to on-hold operation
	existingBook9 := &models.Book{
		ISBN: utils.ToPtr("00009"), 
		State: utils.ToPtr("checked-out"), 
		OnHoldCustomerID: nil,
		CheckedOutCustomerID: utils.ToPtr("08"),
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	daoFactory := &inmemorydao.InMemoryDAOFactory{}
	fixedTimeProvider := &utils.TestingDateTimeProvider{
		ArbitraryTime: arbitraryTimeUpdated,
	}

	h := NewBooksHandler(daoFactory, fixedTimeProvider)
	h.BookDAOInterface.Create(existingBook1)
	h.BookDAOInterface.Create(existingBook2)
	h.BookDAOInterface.Create(existingBook3)
	h.BookDAOInterface.Create(existingBook4)
	h.BookDAOInterface.Create(existingBook5)
	h.BookDAOInterface.Create(existingBook6)
	h.BookDAOInterface.Create(existingBook7)
	h.BookDAOInterface.Create(existingBook8)
	h.BookDAOInterface.Create(existingBook9)
	// h.BookDAOInterface.Create(existingBook10)
	// h.BookDAOInterface.Create(existingBook11)
	// h.BookDAOInterface.Create(existingBook12)
	// h.BookDAOInterface.Create(existingBook13)
	// h.BookDAOInterface.Create(existingBook14)
	// h.BookDAOInterface.Create(existingBook15)
	// h.BookDAOInterface.Create(existingBook16)
	// h.BookDAOInterface.Create(existingBook17)
	// h.BookDAOInterface.Create(existingBook18)
	// h.BookDAOInterface.Create(existingBook19)
	// h.BookDAOInterface.Create(existingBook20)
	// h.BookDAOInterface.Create(existingBook21)
	// h.BookDAOInterface.Create(existingBook22)
	// h.BookDAOInterface.Create(existingBook23)
	// h.BookDAOInterface.Create(existingBook24)
	// h.BookDAOInterface.Create(existingBook25)
	// h.BookDAOInterface.Create(existingBook26)
	// h.BookDAOInterface.Create(existingBook27)
	// h.BookDAOInterface.Create(existingBook28)
	// h.BookDAOInterface.Create(existingBook29)
	// h.BookDAOInterface.Create(existingBook30)


	tests := []struct{
		description string
		currentBook *models.Book
		incomingBook *models.Book
		expectedStatusCode int
		expectedBook *models.Book
		expectedError *models.ErrorResponse
	}{
		{
			description: "Idempotent available to available operation",
			currentBook: existingBook1,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00001"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00001"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: nil,
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: nil,
			},
			expectedError: nil,
		},
		{
			description: "Successfully place hold on an available book",
			currentBook: existingBook2,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00002"),
				State: utils.ToPtr("on-hold"),
				OnHoldCustomerID: utils.ToPtr("04"),
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00002"),
				State: utils.ToPtr("on-hold"),
				OnHoldCustomerID: utils.ToPtr("04"),
				CheckedOutCustomerID: nil,
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: utils.ToPtr(arbitraryTimeUpdated),
			},
			expectedError: nil,
		},
		{
			description: "Successfully checkout an available book",
			currentBook: existingBook3,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00003"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("02"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00003"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("02"),
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: utils.ToPtr(arbitraryTimeUpdated),
			},
			expectedError: nil,
		},
		{
			description: "Successfully return a checked-out book",
			currentBook: existingBook4,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00004"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("02"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00004"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: nil,
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: utils.ToPtr(arbitraryTimeUpdated),
			},
			expectedError: nil,
		},
		{
			description: "Idempotent checked-out to checked-out operation",
			currentBook: existingBook5,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00005"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("02"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00005"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("02"),
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: nil,
			},
			expectedError: nil,
		},
		{
			description: "Successfully release hold on a book",
			currentBook: existingBook6,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00006"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: utils.ToPtr("04"),
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00006"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: nil,
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: utils.ToPtr(arbitraryTimeUpdated),
			},
			expectedError: nil,
		},
		{
			description: "Successfully release hold on a book",
			currentBook: existingBook6,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00006"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: utils.ToPtr("04"),
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00006"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: nil,
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: utils.ToPtr(arbitraryTimeUpdated),
			},
			expectedError: nil,
		},
		{
			description: "Successfully checkout an on-hold book",
			currentBook: existingBook7,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00007"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("06"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00007"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("06"),
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: utils.ToPtr(arbitraryTimeUpdated),
			},
			expectedError: nil,
		},
		{
			description: "Idempotent on-hold to on-hold operation",
			currentBook: existingBook8,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00008"),
				State: utils.ToPtr("on-hold"),
				OnHoldCustomerID: utils.ToPtr("04"),
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 200,
			expectedBook: &models.Book{
				ISBN: utils.ToPtr("00008"),
				State: utils.ToPtr("on-hold"),
				OnHoldCustomerID: utils.ToPtr("04"),
				CheckedOutCustomerID: nil,
				TimeCreated: utils.ToPtr(arbitraryTimeCreated),
				TimeUpdated: nil,
			},
			expectedError: nil,
		},
		{
			description: "Invalid checked-out to on-hold operation",
			currentBook: existingBook9,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00009"),
				State: utils.ToPtr("on-hold"),
				OnHoldCustomerID: utils.ToPtr("08"),
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 409,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Invalid state transition requested: conflict"),
			},
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