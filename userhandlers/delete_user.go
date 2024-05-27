package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteUser allows the client to delete a user by username
func (h *UsersHandler) DeleteUser(c *gin.Context) {
	username := c.Param("username")

	user, err := h.UserDAOInterface.Read(username)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code
		return
	}

	if user == nil {
		c.Status(http.StatusNoContent)
		return
	}

	if err := h.UserDAOInterface.Delete(user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}