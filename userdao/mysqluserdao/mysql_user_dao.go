package mysqluserdao

import (
	// "database/sql"
	"example/library_project/models"
)

type MySQLUserDAO struct {
	// db *sql.DB
}

func (d *MySQLUserDAO) Create(newBook *models.User) error

func (d *MySQLUserDAO) Delete(newBook *models.User) error

func (d *MySQLUserDAO) Update(newBook *models.User) error

func (d *MySQLUserDAO) Read(newBook *models.User) error

func (d *MySQLUserDAO) ReadAll(newBook *models.User) error