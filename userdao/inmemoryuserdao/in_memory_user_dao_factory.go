package inmemoryuserdao

import (
	"example/library_project/models"
	"example/library_project/userdao"
)

type InMemoryUserDAOFactory struct {
	Users map[string]*models.User
}

func NewInMemoryDAOFactory() *InMemoryUserDAOFactory {
	return &InMemoryUserDAOFactory{
		Users: map[string]*models.User{},
	}
}

func (f *InMemoryUserDAOFactory) UserDAO() userdao.UserDAO {
	return &InMemoryUserDAO{
		Users: f.Users,
	}
}

func (f *InMemoryUserDAOFactory) Open() error {
	return nil
}

func (f *InMemoryUserDAOFactory) Close() error {
	return nil
}

func (f *InMemoryUserDAOFactory) Clear() error {
	return nil
}