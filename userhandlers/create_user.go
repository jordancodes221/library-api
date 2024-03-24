package userhandlers

import (
	"github.com/gin-gonic/gin"
)

func (h *UsersHandler) CreateUser(c *gin.Context) {
	// decode JSON to user struct

	// validate

	// hash password, and update user struct from plain text to hashed password
}