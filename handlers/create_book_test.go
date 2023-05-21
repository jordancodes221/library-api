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

func TestBooksHandler_CreateBook(t *testing.T) {
	arbitraryTime := time.Date(2023, 1, 1, 1, 30, 0, 0, time.UTC)

	existingBook := &models.Book{
		ISBN: utils.ToPtr("11111"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil, 
		CheckedOutCustomerID: nil, 
		TimeCreated: utils.ToPtr(arbitraryTime), 
		TimeUpdated: nil,
	}

	daoFactory := &inmemorydao.InMemoryDAOFactory{}
	fixedTimeProvider := &utils.TestingDateTimeProvider{
		ArbitraryTime: arbitraryTime,
	}

	h := NewBooksHandler(daoFactory, fixedTimeProvider)
	h.BookDAOInterface.Create(existingBook)
	
	tests := []struct{
		description string
		book *models.Book
		expectedStatusCode int
		expectedBook *models.Book
		expectedError *models.ErrorResponse
	}{
		{
			description: "valid book", 
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
				TimeCreated: utils.ToPtr(arbitraryTime), 
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
		{
			description: "Invalid state provided",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("invalid state"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Invalid state provided. State must be equal to one of: \"available\", \"on-hold\", or \"checked-out\"."),
			},
		},
		{
			description: "On-hold customer ID is the empty string",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("on-hold"), 
				OnHoldCustomerID: utils.ToPtr(""), 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("On-hold customer ID cannot be the empty string."),
			},
		},
		{
			description: "Checked-out customer ID is the empty string",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("checked-out"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: utils.ToPtr(""), 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Checked-out customer ID cannot be the empty string."),
			},
		},
		{
			description: "Book already exists",
			book: &models.Book{
				ISBN: utils.ToPtr("11111"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 409,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Book already exists."),
			},
		},
		{
			description: "Missing ISBN",
			book: &models.Book{
				ISBN: nil, 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Missing ISBN in the incoming request."),
			},
		},
		{
			description: "Missing State",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: nil, 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Missing State in the incoming request."),
			},
		},
		{
			description: "State is 'available', but on-hold customer ID is non-null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: utils.ToPtr("01"), 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Cannot have an on-hold customer ID when state is available."),
			},
		},
		{
			description: "State is 'available', but checked-out customer ID is non-null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: utils.ToPtr("02"), 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Cannot have checked-out customer ID when state is available."),
			},
		},
		{
			description: "State is 'available', but both on-hold customer ID and checked-out customer ID are non-null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: utils.ToPtr("01"), 
				CheckedOutCustomerID: utils.ToPtr("02"), 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Cannot have an on-hold customer ID when state is available."),
			},
		},
		{
			description: "State is 'on-hold', but checked-out customer ID is non-null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("on-hold"), 
				OnHoldCustomerID: utils.ToPtr("01"), 
				CheckedOutCustomerID: utils.ToPtr("02"), 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Cannot have checked-out customer ID when state is on-hold."),
			},
		},
		{
			description: "State is 'on-hold', but checked-out customer ID is non-null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("on-hold"), 
				OnHoldCustomerID: utils.ToPtr("01"), 
				CheckedOutCustomerID: utils.ToPtr("02"), 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Cannot have checked-out customer ID when state is on-hold."),
			},
		},
		{
			description: "State is 'on-hold', but on-hold customer ID is null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("on-hold"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("State provided is on-hold, but no on-hold customer ID is provided."),
			},
		},
		{
			description: "State is 'on-hold', but checked-out customer ID is non-null and on-hold customer ID is null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("on-hold"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: utils.ToPtr("02"), 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Cannot have checked-out customer ID when state is on-hold."),
			},
		},
		{
			description: "State is 'checked-out', but on-hold customer ID is non-null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("checked-out"), 
				OnHoldCustomerID: utils.ToPtr("01"), 
				CheckedOutCustomerID: utils.ToPtr("02"), 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Cannot have on-hold customer ID when state is checked-out."),
			},
		},
		{
			description: "State is 'checked-out', but checked-out customer ID is null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("checked-out"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("State provided is checked-out, but no checked-out customer ID is provided."),
			},
		},
		{
			description: "State is 'checked-out', but on-hold customer ID is non-null and checked-out customer ID is null",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("checked-out"), 
				OnHoldCustomerID: utils.ToPtr("01"), 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Cannot have on-hold customer ID when state is checked-out."),
			},
		},
		{
			description: "Time Created is provided",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: utils.ToPtr(time.Now()), 
				TimeUpdated: nil,
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Client cannot provide time created when creating a new book."),
			},
		},
		{
			description: "Time Updated is provided",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: nil, 
				TimeUpdated: utils.ToPtr(time.Now()),
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Client cannot provide time updated when creating a new book."),
			},
		},
		{
			description: "Time Created and Time Updated are both provided",
			book: &models.Book{
				ISBN: utils.ToPtr("00000"), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: utils.ToPtr(time.Now()), 
				TimeUpdated: utils.ToPtr(time.Now()),
			}, 
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Client cannot provide time created when creating a new book."),
			},
		},

	}
	
	for _, currentTestCase := range tests {
		fmt.Println(currentTestCase.description)
		t.Log(currentTestCase.description)

		bookJSON, _ := json.Marshal(*currentTestCase.book)

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
