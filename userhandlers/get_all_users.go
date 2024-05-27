package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllUsers allows the client to get all of the users
func (h *UsersHandler) GetAllUsers(c *gin.Context) {
	all_users, err := h.UserDAOInterface.ReadAll()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()}) // 500 status code if unsuccessful
		return
	}

	c.IndentedJSON(http.StatusOK, all_users)
}