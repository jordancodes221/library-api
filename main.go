package main

import ( 
	"example/library_project/handlers"
	"example/library_project/models"
	"example/library_project/utils"
	
	// "net/http"
	"github.com/gin-gonic/gin"
	// "errors"
	"time"
	// "encoding/json"
	// "fmt"

	// "reflect"
	// "strconv"
)

func main() {

	//// Instantiating Test Data

	// First test of instantiating test data with new schema and utils.ToPtr function
	var bookInstance00 *models.Book = &models.Book{ISBN: utils.ToPtr("00"), State: utils.ToPtr("on-hold"), OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Now())}

	// Actual test data to be used in testing
	var bookInstance0 *models.Book = &models.Book{ISBN: utils.ToPtr("0000"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Available
	var bookInstance1 *models.Book = &models.Book{ISBN: utils.ToPtr("0001"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Checked-out
	var bookInstance2 *models.Book = &models.Book{ISBN: utils.ToPtr("0002"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> On-hold

	var bookInstance3 *models.Book = &models.Book{ISBN: utils.ToPtr("0003"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Available
	var bookInstance4 *models.Book = &models.Book{ISBN: utils.ToPtr("0004"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Available (no match)
	var bookInstance5 *models.Book = &models.Book{ISBN: utils.ToPtr("0005"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Checked-out
	var bookInstance6 *models.Book = &models.Book{ISBN: utils.ToPtr("0006"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Checked-out (no match)
	var bookInstance7 *models.Book = &models.Book{ISBN: utils.ToPtr("0007"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> On-hold 
	var bookInstance8 *models.Book = &models.Book{ISBN: utils.ToPtr("0008"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> On-hold (no match)

	var bookInstance9 *models.Book =  &models.Book{ISBN: utils.ToPtr("0009"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Available
	var bookInstance10 *models.Book = &models.Book{ISBN: utils.ToPtr("0010"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Available (no match)
	var bookInstance11 *models.Book = &models.Book{ISBN: utils.ToPtr("0011"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Checked-out
	var bookInstance12 *models.Book = &models.Book{ISBN: utils.ToPtr("0012"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> Checked-out (no match)
	var bookInstance13 *models.Book = &models.Book{ISBN: utils.ToPtr("0013"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> On-hold 
	var bookInstance14 *models.Book = &models.Book{ISBN: utils.ToPtr("0014"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> On-hold (no match)

	var bookInstance15 *models.Book = &models.Book{ISBN: utils.ToPtr("0015"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, 	TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} // --> This is the book to be deleted in testing

	// The following are for UpdateBook ID semantics validation
	var bookInstance16 *models.Book = &models.Book{ISBN: utils.ToPtr("0016"), State: utils.ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, 	TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})} 
	var bookInstance17 *models.Book = &models.Book{ISBN: utils.ToPtr("0017"), State: utils.ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: utils.ToPtr("01"), TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})}
	var bookInstance18 *models.Book = &models.Book{ISBN: utils.ToPtr("0018"), State: utils.ToPtr("on-hold"), 	OnHoldCustomerID: utils.ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Now()), TimeUpdated: utils.ToPtr(time.Time{})}

	// The follwing are for UpdateBook Time validation
	arbitraryIncomingTimeCreated, _ := time.Parse(time.RFC3339, "2023-03-18T15:45:00Z")
	arbitraryIncomingTimeUpdated, _ := time.Parse(time.RFC3339, "2022-02-18T15:45:00Z")
	// Notes:
		// (1) The body of the requests in Postman all send the above time created and updated.
		// (2) The test data below has been instantiated with select time field set to zero (via time.Time{}) to intentionally create a mismatch for our testing.
	var bookInstance19 *models.Book = &models.Book{ISBN: utils.ToPtr("0019"), State: utils.ToPtr("available"), 	OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Time{}), TimeUpdated: utils.ToPtr(arbitraryIncomingTimeUpdated)}
	var bookInstance20 *models.Book = &models.Book{ISBN: utils.ToPtr("0020"), State: utils.ToPtr("available"), 	OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(arbitraryIncomingTimeCreated), TimeUpdated: utils.ToPtr(time.Time{})}
	var bookInstance21 *models.Book = &models.Book{ISBN: utils.ToPtr("0021"), State: utils.ToPtr("available"), 	OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(arbitraryIncomingTimeCreated), TimeUpdated: utils.ToPtr(arbitraryIncomingTimeUpdated)}
	var bookInstance22 *models.Book = &models.Book{ISBN: utils.ToPtr("0022"), State: utils.ToPtr("available"), 	OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: utils.ToPtr(time.Time{}), TimeUpdated: utils.ToPtr(time.Time{})}

	// Map of test data to be used in testing
	var mapOfBooks = map[string]*models.Book{
		"00" : bookInstance00,

		"0000" : bookInstance0,
		"0001" : bookInstance1,
		"0002" : bookInstance2,

		"0003" : bookInstance3,
		"0004" : bookInstance4,
		"0005" : bookInstance5,
		"0006" : bookInstance6,
		"0007" : bookInstance7,
		"0008" : bookInstance8,

		"0009" : bookInstance9,
		"0010" : bookInstance10,
		"0011" : bookInstance11,
		"0012" : bookInstance12,
		"0013" : bookInstance13,
		"0014" : bookInstance14,

		"0015" : bookInstance15,

		"0016" : bookInstance16,
		"0017" : bookInstance17,
		"0018" : bookInstance18,

		"0019" : bookInstance19,
		"0020" : bookInstance20,
		"0021" : bookInstance21,	
		"0022" : bookInstance22,		
	}

	router := gin.Default()
	h := &handlers.BooksHandler{Books: mapOfBooks}
	router.GET("/books", h.GetAllBooks)
	router.GET("/books/:isbn", h.GetIndividualBook)
	router.POST("/books", h.CreateBook)
	router.DELETE("/books/:isbn", h.DeleteBook)
	router.PATCH("/books/:isbn", h.UpdateBook)

	router.Run("localhost:8080")
}

// To test, run "go run ." in one terminal window and a curl command in the another terminal window.
// Examples of curl commands are:
	// GET (all books)
		// curl localhost:8080/books
	// GET (individual book)
		// curl localhost:8080/books/0000
	// POST
		// curl localhost:8080/books --include --header "Content-Type: application/json" -d @newBook.json --request "POST"
	// DELETE
		// curl localhost:8080/books/0005 --request "DELETE"
	// PATCH
		// curl -X PATCH localhost:8080/books/00 -H 'Content-Type: application/json' -H 'Accept: application/json' -d @incomingRequest.json