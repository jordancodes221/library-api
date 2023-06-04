package mysqldao


import (
	"database/sql"
	"example/library_project/models"

	// "fmt"
	// "log"
)

type MySQLBookDAO struct {
	db *sql.DB
}

func (d *MySQLBookDAO) Create(newBook *models.Book) {
	return
}

func (d *MySQLBookDAO) Delete(newBook *models.Book) {
	return
}

func (d *MySQLBookDAO) Update(newBook *models.Book) {
	return
}

func (d *MySQLBookDAO) Read(isbn string) (*models.Book, error) {
	return nil, nil
}

func (d *MySQLBookDAO) ReadAll() ([]*models.Book, error) {
	return nil, nil
}