package mysqldao


import (
	"database/sql"
	"example/library_project/models"

	"fmt"
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
	query := "SELECT ISBN, State, OnHoldCustomerID, CheckedOutCustomerID, TimeCreated, TimeUpdated FROM Books"

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	retrievedBooks := make([]*models.Book, 0)

	for rows.Next() {
		nextBook := new(models.Book)

		err := rows.Scan(
			&nextBook.ISBN,
			&nextBook.State,
			&nextBook.OnHoldCustomerID,
			&nextBook.CheckedOutCustomerID,
			&nextBook.TimeCreated,
			&nextBook.TimeUpdated,)
		if err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		retrievedBooks = append(retrievedBooks, nextBook)
	}

	return retrievedBooks, nil
}