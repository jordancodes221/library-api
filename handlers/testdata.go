package handlers

import ( 
	// "example/library_project/handlers"
	// "example/library_project/validators"
	"example/library_project/models"
	
	// "net/http"
	// "github.com/gin-gonic/gin"
	// "errors"
	"time"
	// "encoding/json"
	// "fmt"
	// "reflect"
	// "strconv"
)

//// Instantiating Test Data

// First test of instantiating test data with new schema and ToPtr function
var bookInstance00 models.Book = models.Book{ISBN: ToPtr("00"), State: ToPtr("on-hold"), OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Now())}

// Actual test data to be used in testing
var bookInstance0 models.Book = models.Book{ISBN: ToPtr("0000"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available
var bookInstance1 models.Book = models.Book{ISBN: ToPtr("0001"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out
var bookInstance2 models.Book = models.Book{ISBN: ToPtr("0002"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold

var bookInstance3 models.Book = models.Book{ISBN: ToPtr("0003"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available
var bookInstance4 models.Book = models.Book{ISBN: ToPtr("0004"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available (no match)
var bookInstance5 models.Book = models.Book{ISBN: ToPtr("0005"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out
var bookInstance6 models.Book = models.Book{ISBN: ToPtr("0006"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out (no match)
var bookInstance7 models.Book = models.Book{ISBN: ToPtr("0007"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold 
var bookInstance8 models.Book = models.Book{ISBN: ToPtr("0008"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold (no match)

var bookInstance9 models.Book =  models.Book{ISBN: ToPtr("0009"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available
var bookInstance10 models.Book = models.Book{ISBN: ToPtr("0010"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Available (no match)
var bookInstance11 models.Book = models.Book{ISBN: ToPtr("0011"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out
var bookInstance12 models.Book = models.Book{ISBN: ToPtr("0012"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> Checked-out (no match)
var bookInstance13 models.Book = models.Book{ISBN: ToPtr("0013"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold 
var bookInstance14 models.Book = models.Book{ISBN: ToPtr("0014"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> On-hold (no match)

var bookInstance15 models.Book = models.Book{ISBN: ToPtr("0015"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, 	TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} // --> This is the book to be deleted in testing

// The following are for UpdateBook ID semantics validation
var bookInstance16 models.Book = models.Book{ISBN: ToPtr("0016"), State: ToPtr("available"), OnHoldCustomerID: nil, CheckedOutCustomerID: nil, 	TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})} 
var bookInstance17 models.Book = models.Book{ISBN: ToPtr("0017"), State: ToPtr("checked-out"), OnHoldCustomerID: nil, CheckedOutCustomerID: ToPtr("01"), TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})}
var bookInstance18 models.Book = models.Book{ISBN: ToPtr("0018"), State: ToPtr("on-hold"), 	OnHoldCustomerID: ToPtr("01"), CheckedOutCustomerID: nil, TimeCreated: ToPtr(time.Now()), TimeUpdated: ToPtr(time.Time{})}

// Map of test data to be used in testing
var mapOfBooks = map[string]*models.Book{
	"00" : &bookInstance00,

	"0000" : &bookInstance0,
	"0001" : &bookInstance1,
	"0002" : &bookInstance2,

	"0003" : &bookInstance3,
	"0004" : &bookInstance4,
	"0005" : &bookInstance5,
	"0006" : &bookInstance6,
	"0007" : &bookInstance7,
	"0008" : &bookInstance8,

	"0009" : &bookInstance9,
	"0010" : &bookInstance10,
	"0011" : &bookInstance11,
	"0012" : &bookInstance12,
	"0013" : &bookInstance13,
	"0014" : &bookInstance14,

	"0015" : &bookInstance15,

	"0016" : &bookInstance16,
	"0017" : &bookInstance17,
	"0018" : &bookInstance18,
}