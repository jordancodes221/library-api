# Library API

This REST API handles the state of books in a library. It is written in Go and utilizes the Gin framework. MySQL is implemented as one possible storage solution.

## Design

- The handlers package contains handler functions that implement the HTTP methods GET, PUT, POST, and DELETE.
  - Notably the UpdateBook handler function does not simply toggle individual fields of the book resource. Instead, it compares the requested state to the current state to determine whether to update the current state to the requested one.
  - Each handler function has associated validator functions that perform syntax and logic validation.
- The data access object (DAO) contains the create, read, update and delete (CRUD) functions that interact with the storage layer.
  - This abstraction of the CRUD functions from the handler functions eased scalabiilty as the handler functions do not need to be re-written when the storage solution is changed (such as when I scaled the API from in-memory to MySQL storage).
  - The Abstract Factory design pattern was followed for implementing different versions of the DAO for different storage solutions.

## Testing

- Unit tests are implemented for each handler function, as well as for the book model. These can be found as separate test files in the respective package.
- Postman was used to perform integration tests.
