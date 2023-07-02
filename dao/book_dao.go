package dao

import (
	"example/library_project/models"
)

type BookDAO interface {
	// once a persistent database is added, these methods will also return an error type
	Create(newBook *models.Book) error
	Read(isbn string) (*models.Book, error)
	ReadAll() ([]*models.Book, error)
	Update(book *models.Book) error
	Delete(book *models.Book) error
}