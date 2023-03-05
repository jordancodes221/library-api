package handlers

import (
	"example/library_project/models"
)

type BooksHandler struct {
    Books map[string]*models.Book
}

