package userhandlers

import (
	"errors"
	"example/library_project/models"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func validateLogicForUpdateUser(incomingUser *models.User, currentUser *models.User) (error) {	
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

func (h *UsersHandler) UpdateUser(c *gin.Context) { 
	username := c.Param("username")

	currentUser, err := h.UserDAOInterface.Read(username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if currentUser == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "User not found."})
		return
	}	

	// Decode JSON to user struct
	incomingUser := new(models.User) // the "new" keyword allocates memory for models.User, and returns a pointer to it
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(incomingUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// If fields are not nil, ensure they are within range
	if err := incomingUser.Validate(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Validate logic
	if err := validateLogicForUpdateUser(incomingUser, currentUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	// Update the user
	if err := h.UserDAOInterface.Update(currentUser); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, currentUser)
}