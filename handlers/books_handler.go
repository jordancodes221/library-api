package handlers

import (
	"example/library_project/models"
	"example/library_project/utils"
	"example/library_project/dao"
)

// BooksHandlers is the struct on which all handler functions are defined as pointer-receiver functions
type BooksHandler struct {
	// Books is the library of all the books
    Books map[string]*models.Book
	DateTimeInterface utils.DateTimeProvider
}