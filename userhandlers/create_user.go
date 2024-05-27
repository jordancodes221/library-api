package userhandlers

import (
	"encoding/json"
	"example/library_project/models"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

func validateLogicForCreateUser(incomingUser *models.User) (error) {
	// Ensure username is provided
	if incomingUser.Username == nil {
		return errors.New("missing username in the incoming request")
	}
	
	// Ensure password is provided
	if incomingUser.Password == nil {
		return errors.New("missing password in the incoming request")
	}

	// Verify password is at least 8 characters long
	if len(*incomingUser.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Could verify other conditions (i.e. contains special character, etc)

	return nil
}

func (h *UsersHandler) CreateUser(c *gin.Context) {
	// decode JSON to user struct
	newUser := new(models.User) // the "new" keyword allocates memory for models.User, and returns a pointer to it
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// If fields are not nil, ensure they are within range
	if err := newUser.Validate(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Logic validation
	if err := validateLogicForCreateUser(newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Ensure username is not already in use
	bookWithUsernameInUse, err := h.UserDAOInterface.Read(*newUser.Username)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if bookWithUsernameInUse != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"ERROR": "username already exists"})
		return
	}

	// TO DO: hash password, and update user struct from plain text to hashed password

	// Add the new book to our library
	if err := h.UserDAOInterface.Create(newUser); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser) // 201 status code if successful
}