package mysqldao


import (
	"database/sql"
	"example/library_project/models"

	"fmt"
	// "log"

	"time"
)

type MySQLBookDAO struct {
	db *sql.DB
}

func (d *MySQLBookDAO) Create(newBook *models.Book) {
	query := "INSERT INTO Books (ISBN, State, OnHoldCustomerID, CheckedOutCustomerID, TimeCreated, TimeUpdated) VALUES (?, ?, ?, ?, ?, ?)"

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

func (d *MySQLBookDAO) Update(book *models.Book) {
	query := "UPDATE Books SET State = ?, OnHoldCustomerID = ?, CheckedOutCustomerID = ?, TimeUpdated = ? WHERE ISBN = ?"

	_, err := d.db.Exec(query, book.State, book.OnHoldCustomerID, book.CheckedOutCustomerID, book.TimeUpdated, book.ISBN)
	if err != nil {
		fmt.Println("error updating book: ", err)
		return
	}

	return
}

func (d *MySQLBookDAO) Read(isbn string) (*models.Book, error) {
	query := "SELECT ISBN, State, OnHoldCustomerID, CheckedOutCustomerID, TimeCreated, TimeUpdated FROM Books WHERE ISBN = ?"
	row := d.db.QueryRow(query, isbn)

	retrievedISBN := new(sql.NullString)
	retrievedState := new(sql.NullString)
	retrievedOnHoldCustomerID := new(sql.NullString)
	retrievedCheckedOutCustomerID := new(sql.NullString)
	retrievedTimeCreated := new(sql.NullString)
	retrievedTimeUpdated := new(sql.NullString)

	err := row.Scan(
		retrievedISBN,
		retrievedState,
		retrievedOnHoldCustomerID,
		retrievedCheckedOutCustomerID,
		retrievedTimeCreated,
		retrievedTimeUpdated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Book not found: %v", err)
		}
		
		return nil, fmt.Errorf("error: %v", err)
	}

	retrievedIndividualBook := &models.Book{
		ISBN: nil,
		State: nil,
		OnHoldCustomerID: nil,
		CheckedOutCustomerID: nil,
		TimeCreated: nil,
		TimeUpdated: nil,
	}

	if retrievedISBN.Valid {
		retrievedIndividualBook.ISBN = &retrievedISBN.String
	}

	if retrievedState.Valid {
		retrievedIndividualBook.State = &retrievedState.String
	}

	if retrievedOnHoldCustomerID.Valid {
		retrievedIndividualBook.OnHoldCustomerID = &retrievedOnHoldCustomerID.String
	}

	if retrievedCheckedOutCustomerID.Valid {
		retrievedIndividualBook.CheckedOutCustomerID = &retrievedCheckedOutCustomerID.String
	}

	if retrievedTimeCreated.Valid {
		timeCreated, err := time.Parse("2006-01-02 15:04:05", retrievedTimeCreated.String)
		if err != nil {
			return nil, fmt.Errorf("error parsing time created in read: %w", err)
		}
		retrievedIndividualBook.TimeCreated = &timeCreated
	}

	if retrievedTimeUpdated.Valid {
		timeUpdated, err := time.Parse("2006-01-02 15:04:05", retrievedTimeUpdated.String)
		if err != nil {
			return nil, fmt.Errorf("error parsing time created in read: %w", err)
		}
		retrievedIndividualBook.TimeUpdated = &timeUpdated
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