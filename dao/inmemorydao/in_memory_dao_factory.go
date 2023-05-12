package inmemorydao

import (
	"example/library_project/dao"
	"example/library_project/models"
)

type InMemoryDAOFactory struct {}

func (f *InMemoryDAOFactory) BookDAO() dao.BookDAO {
	return &InMemoryBookDAO{
		Books: make(map[string]*models.Book),
	}
}