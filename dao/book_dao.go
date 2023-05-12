package dao

import (
	"example/library_project/models"
)

type BookDAO interface {
	// once a persistent database is added, these methods will also return an error type
	Create(newBook *models.Book)
	Read(isbn string) (*models.Book, error)
	ReadAll() ([]*models.Book, error)
	Update(book *models.Book)
	Delete(book *models.Book)
}