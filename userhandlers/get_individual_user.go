package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndividualUser allows the client to get an individual user by username
func (h *UsersHandler) GetIndividualUser(c *gin.Context) {
	username := c.Param("username")
	user, err := h.UserDAOInterface.Read(username)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code if unsuccessful
		return
	}

	if user == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "REQUEST SUCCESSFUL. BOOK NOT FOUND"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}