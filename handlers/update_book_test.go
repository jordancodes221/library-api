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

	existingBook10 := &models.Book{
		ISBN: utils.ToPtr("000010"), 
		State: utils.ToPtr("checked-out"), 
		OnHoldCustomerID: nil,
		CheckedOutCustomerID: utils.ToPtr("10"),
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook11 := &models.Book{
		ISBN: utils.ToPtr("000011"), 
		State: utils.ToPtr("checked-out"), 
		OnHoldCustomerID: nil,
		CheckedOutCustomerID: utils.ToPtr("10"),
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook12 := &models.Book{
		ISBN: utils.ToPtr("000012"), 
		State: utils.ToPtr("checked-out"), 
		OnHoldCustomerID: nil,
		CheckedOutCustomerID: utils.ToPtr("10"),
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook13 := &models.Book{
		ISBN: utils.ToPtr("000013"), 
		State: utils.ToPtr("on-hold"), 
		OnHoldCustomerID: utils.ToPtr("10"),
		CheckedOutCustomerID: nil,
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook14 := &models.Book{
		ISBN: utils.ToPtr("000014"), 
		State: utils.ToPtr("on-hold"), 
		OnHoldCustomerID: utils.ToPtr("10"),
		CheckedOutCustomerID: nil,
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook15 := &models.Book{
		ISBN: utils.ToPtr("000015"), 
		State: utils.ToPtr("on-hold"), 
		OnHoldCustomerID: utils.ToPtr("10"),
		CheckedOutCustomerID: nil,
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook16 := &models.Book{
		ISBN: utils.ToPtr("000016"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil,
		CheckedOutCustomerID: nil,
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	existingBook17 := &models.Book{
		ISBN: utils.ToPtr("000017"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil,
		CheckedOutCustomerID: nil,
		TimeCreated: utils.ToPtr(arbitraryTimeCreated), 
		TimeUpdated: nil,
	}

	// existingBook18 is used for the "Book not found" test case so do not add it do the BookDAO
	existingBook18 := &models.Book{
		ISBN: utils.ToPtr("000018"), 
		State: utils.ToPtr("available"), 
		OnHoldCustomerID: nil,
		CheckedOutCustomerID: nil,
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
	h.BookDAOInterface.Create(existingBook10)
	h.BookDAOInterface.Create(existingBook11)
	h.BookDAOInterface.Create(existingBook12)
	h.BookDAOInterface.Create(existingBook13)
	h.BookDAOInterface.Create(existingBook14)
	h.BookDAOInterface.Create(existingBook15)
	h.BookDAOInterface.Create(existingBook16)
	h.BookDAOInterface.Create(existingBook17)
	// existingBook18 is used for the "Book not found" test case
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
			description: "Invalid checked-out to on-hold operation (invalid state transition)",
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
		{
			description: "Invalid checked-out to available (IDs do not match)",
			currentBook: existingBook10,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00010"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("20"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 409,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Returning the book failed as it is another customer who has the book checked-out: conflict"),
			},
		},
		{
			description: "Invalid checked-out to checked-out (IDs do not match)",
			currentBook: existingBook11,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00011"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("20"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 409,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Checkout failed as another customer has the book checked-out: conflict"),
			},
		},
		{
			description: "Invalid checked-out to on-hold (invalid state transition and IDs do not match)",
			currentBook: existingBook12,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00012"),
				State: utils.ToPtr("on-hold"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("20"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 409,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Invalid state transition requested: conflict"),
			},
		},
		{
			description: "Invalid on-hold to available (IDs do not match)",
			currentBook: existingBook13,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00013"),
				State: utils.ToPtr("available"),
				OnHoldCustomerID: utils.ToPtr("20"),
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 409,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Releasing hold failed as it is another customer who has the book on-hold: conflict"),
			},
		},
		{
			description: "Invalid on-hold to checked-out (IDs do not match)",
			currentBook: existingBook14,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00014"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("20"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 409,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Checkout failed as another customer has the book on-hold: conflict"),
			},
		},
		{
			description: "Invalid on-hold to on-hold (IDs do not match)",
			currentBook: existingBook15,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00015"),
				State: utils.ToPtr("on-hold"),
				OnHoldCustomerID: utils.ToPtr("20"),
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 409,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Placing hold failed as another customer has the book on-hold: conflict"),
			},
		},
		{
			description: "Missing state in request",
			currentBook: existingBook16,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00016"),
				State: nil,
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: nil,
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 400,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Expected 'state' to be non-null: invalid request"),
			},
		},
		{
			description: "Invalid state in request",
			currentBook: existingBook17,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00017"),
				State: utils.ToPtr("invalid-state"),
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
			description: "Book not found",
			currentBook: existingBook18,
			incomingBook: &models.Book{
				ISBN: utils.ToPtr("00018"),
				State: utils.ToPtr("checked-out"),
				OnHoldCustomerID: nil,
				CheckedOutCustomerID: utils.ToPtr("100"),
				TimeCreated: nil,
				TimeUpdated: nil,
			},
			expectedStatusCode: 404,
			expectedBook: nil,
			expectedError: &models.ErrorResponse{
				Message: utils.ToPtr("Book not found."),
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