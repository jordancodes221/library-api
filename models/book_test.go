package models

import (
	"testing"
	"time"
	"example/library_project/utils"
	"github.com/stretchr/testify/assert"
)

func TestBook_Validate(t *testing.T){
	tests := []struct{
		description string
		book *Book
		expectedErrorMessage string
	}{
		{
			description: "ISBN is the empty string", 
			book: &Book{
				ISBN: utils.ToPtr(""), 
				State: utils.ToPtr("available"), 
				OnHoldCustomerID: nil, 
				CheckedOutCustomerID: nil, 
				TimeCreated: utils.ToPtr(time.Now()), 
				TimeUpdated: utils.ToPtr(time.Time{}),
			}, 
			expectedErrorMessage: "ISBN cannot be the empty string.",
		},
	}

	for _, currentTestCase := range tests {
		t.Log(currentTestCase.description)
		actual := currentTestCase.book.Validate()
		assert.EqualError(t, actual, currentTestCase.expectedErrorMessage)
	}
}