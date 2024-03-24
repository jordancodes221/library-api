package userdao

import (
	"example/library_project/models"
)

type UserDAO interface {
	// Create adds a new user to the library
	Create(newUser *models.User) error
	Read(isbn string) (*models.User, error)
	ReadAll() ([]*models.User, error)
	Update(book *models.User) error
	Delete(book *models.User) error
}