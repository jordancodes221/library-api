package inmemorydao

import (
	"example/library_project/dao"
	"example/library_project/models"
)

type InMemoryDAOFactory struct {
	Books map[string]*models.Book
}

func NewInMemoryDAOFactory() *InMemoryDAOFactory {
	return &InMemoryDAOFactory{
		Books: map[string]*models.Book{},
	}
}

func (f *InMemoryDAOFactory) BookDAO() dao.BookDAO {
	return &InMemoryBookDAO{
		Books: f.Books,
	}
}

func (f *InMemoryDAOFactory) Open() error {
	return nil
}

func (f *InMemoryDAOFactory) Close() error {
	return nil
}

func (f *InMemoryDAOFactory) Clear() error {
	return nil
}