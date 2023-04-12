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
		// expectedBody interface{}
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
			// expectedBody: &models.Book{
			// 	ISBN: utils.ToPtr("00000"), 
			// 	State: utils.ToPtr("available"), 
			// 	OnHoldCustomerID: nil, 
			// 	CheckedOutCustomerID: nil, 
			// 	TimeCreated: nil, 
			// 	TimeUpdated: nil,
			// },
		}, 
		// {
		// 	description: "ISBN is the empty string", 
		// 	book: &models.Book{
		// 		ISBN: utils.ToPtr(""), 
		// 		State: utils.ToPtr("available"), 
		// 		OnHoldCustomerID: nil, 
		// 		CheckedOutCustomerID: nil, 
		// 		TimeCreated: nil, 
		// 		TimeUpdated: nil,
		// 	}, 
		// 	expectedStatusCode: 400,
		// 	// expectedBody: "ISBN cannot be the empty string.",
		// },
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
		// assert.Equal(t, currentTestCase.expectedBody, w.Body)
	}
}
