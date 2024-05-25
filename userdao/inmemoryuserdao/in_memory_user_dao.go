package inmemoryuserdao

import (
	"example/library_project/models"
)

type InMemoryUserDAO struct {
	Users map[string]*models.User
}

func (d *InMemoryUserDAO) Create(newUser *models.User) error {
	d.Users[*newUser.Userid] = newUser
	return nil
}

func (d *InMemoryUserDAO) Delete(book *models.User) error {
	delete(d.Users, *book.Userid)
	return nil
}

func (d *InMemoryUserDAO) Update(book *models.User) error {
	d.Users[*book.Userid] = book
	return nil
}

func (d *InMemoryUserDAO) Read(userid string) (*models.User, error) {
	retrievedUser, ok := d.Users[userid] // in the future, this could be a call to a database

	// For scalability, we can add a database connection here. 
	// If there is an error connecting to the database, then we will return: nil, InternalServerError

	if ok {
		return retrievedUser, nil
	} else {
		return nil, nil
	}
}

func (d *InMemoryUserDAO) ReadAll() ([]*models.User, error) {	
	all_users := make([]*models.User, 0)

	// For scalability, we can add a database connection here. 
	// If there is an error connecting to the database, then we will return: nil, InternalServerError

	for _, currentUser := range d.Users {
		all_users = append(all_users, currentUser)
	}

	return all_users, nil
}