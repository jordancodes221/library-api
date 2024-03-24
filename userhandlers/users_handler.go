package userhandlers

import (
	"example/library_project/userdao"
	"example/library_project/utils"
)

// UsersHandlers is the struct on which all handler functions are defined as pointer-receiver functions
type UsersHandler struct {
	// Users is the library of all the users
	UserDAOInterface userdao.UserDAO
	DateTimeInterface utils.DateTimeProvider
}

func NewUsersHandler(userDAO userdao.UserDAO, provider utils.DateTimeProvider) (*UsersHandler) {
	return &UsersHandler{
		UserDAOInterface: userDAO,
		DateTimeInterface: provider,
	}
}
