package testdata

import (
	"example/library_project/models"
	"example/library_project/utils"
	"time"
)

func InstantiateIntegrationTestData() ([]*models.Book, error) {

	// The following are for the instances used for testing UpdateBook Time validation
	arbitraryIncomingTimeCreated, _ := time.Parse(time.RFC3339, "2023-03-18T15:45:00Z")
	arbitraryIncomingTimeUpdated, _ := time.Parse(time.RFC3339, "2022-02-18T15:45:00Z")

	integrationTestData := []*models.Book{
		// First test of instantiating test data with new schema and utils.ToPtr function
		{ISBN: utils.ToPtr("00"), State: utils.ToPtr("on-hold"), OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Now())}, 

		// Actual test data to be used in testing
		{ISBN: utils.ToPtr("0000"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Available
		{ISBN: utils.ToPtr("0001"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Checked-out
		{ISBN: utils.ToPtr("0002"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> On-hold

		{ISBN: utils.ToPtr("0003"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Available
		{ISBN: utils.ToPtr("0004"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Available (no match)
		{ISBN: utils.ToPtr("0005"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Checked-out
		{ISBN: utils.ToPtr("0006"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Checked-out (no match)
		{ISBN: utils.ToPtr("0007"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> On-hold 
		{ISBN: utils.ToPtr("0008"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> On-hold (no match)

		{ISBN: utils.ToPtr("0009"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Available
		{ISBN: utils.ToPtr("0010"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Available (no match)
		{ISBN: utils.ToPtr("0011"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Checked-out
		{ISBN: utils.ToPtr("0012"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> Checked-out (no match)
		{ISBN: utils.ToPtr("0013"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> On-hold 
		{ISBN: utils.ToPtr("0014"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  // --> On-hold (no match)

		// The following instance is the book to be deleted when testing DeleteBook
		{ISBN: utils.ToPtr("0015"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, 	TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil}, 

		// The following are for UpdateBook ID semantics validation
		{ISBN: utils.ToPtr("0016"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, 	TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil},  
		{ISBN: utils.ToPtr("0017"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil}, 
		{ISBN: utils.ToPtr("0018"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: nil}, 

		// Notes:
			// (1) The body of the requests in Postman all send the above time created and updated.
			// (2) The test data below has been instantiated with select time field set to zero (via time.Time{}) to intentionally create a mismatch for our testing.
		{ISBN: utils.ToPtr("0019"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Time{}), TimeUpdated: utils.ToPtr(arbitraryIncomingTimeUpdated)}, 
		{ISBN: utils.ToPtr("0020"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(arbitraryIncomingTimeCreated), TimeUpdated: utils.ToPtr(time.Time{})}, 
		{ISBN: utils.ToPtr("0021"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(arbitraryIncomingTimeCreated), TimeUpdated: utils.ToPtr(arbitraryIncomingTimeUpdated)}, 
		{ISBN: utils.ToPtr("0022"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Time{}), TimeUpdated: utils.ToPtr(time.Time{})}, 

	}

	return integrationTestData, nil
}