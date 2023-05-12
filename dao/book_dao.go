package dao

import (
	"example/library_project/models"
)

type BookDAO interface {
	// once a persistent database is added, these methods will also return an error type
	Create()
	Read() *models.Book
	ReadAll() []*models.Book
	Update()
	Delete()
}