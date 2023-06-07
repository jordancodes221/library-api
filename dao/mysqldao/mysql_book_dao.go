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
	query := "INSERT INTO Books VALUES (?, ?, ?, ?, ?, ?)"

	_, err := d.db.Exec(query, newBook.ISBN, newBook.State, newBook.OnHoldCustomerID, newBook.CheckedOutCustomerID, newBook.TimeCreated, newBook.TimeUpdated)
	if err != nil {
		fmt.Println("error adding new book to database: ", err)
		return
	}

	return
}

func (d *MySQLBookDAO) Delete(book *models.Book) {
	query := "DELETE FROM Books WHERE ISBN = ?"

	_, err := d.db.Exec(query, book.ISBN)
	if err != nil {
		fmt.Println("error deleting book from database: ", err)
		return
	}

	return
}

func (d *MySQLBookDAO) Update(newBook *models.Book) {
	return
}

func (d *MySQLBookDAO) Read(isbn string) (*models.Book, error) {
	query := "SELECT ISBN, State, OnHoldCustomerID, CheckedOutCustomerID, TimeCreated, TimeUpdated FROM Books WHERE ISBN = ?"
	row := d.db.QueryRow(query, isbn)

	retrievedIndividualBook := new(models.Book)

	err := row.Scan(
		&retrievedIndividualBook.ISBN,
		&retrievedIndividualBook.State,
		&retrievedIndividualBook.OnHoldCustomerID,
		&retrievedIndividualBook.CheckedOutCustomerID,
		&retrievedIndividualBook.TimeCreated,
		&retrievedIndividualBook.TimeUpdated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Book not found: %v", err)
		}
		
		return nil, fmt.Errorf("error: %v", err)
	}

	return retrievedIndividualBook, nil
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