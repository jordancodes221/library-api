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

func (d *MySQLBookDAO) Create(newBook *models.Book) error {
	// query := "INSERT INTO Books (ISBN, State, OnHoldCustomerID, CheckedOutCustomerID, TimeCreated, TimeUpdated) VALUES (?, ?, ?, ?, ?, ?)"
	query := "INSERT INTO Books (ISBN, State, OnHoldCustomerID, CheckedOutCustomerID, TimeCreated, TimeUpdated) VALUES (?, ?, ?, ?, ?, NULL)"

	timeCreated := newBook.TimeCreated.Format("2006-01-02 15:04:05")
	// timeUpdated := newBook.TimeUpdated.Format("2006-01-02 15:04:05")

	// _, err := d.db.Exec(query, newBook.ISBN, newBook.State, newBook.OnHoldCustomerID, newBook.CheckedOutCustomerID, timeCreated, timeUpdated)
	_, err := d.db.Exec(query, newBook.ISBN, newBook.State, newBook.OnHoldCustomerID, newBook.CheckedOutCustomerID, timeCreated)
	if err != nil {
		return fmt.Errorf("error adding new book to database: %w", err)
	}

	return nil
}

func (d *MySQLBookDAO) Delete(book *models.Book) error {
	query := "DELETE FROM Books WHERE ISBN = ?"

	_, err := d.db.Exec(query, book.ISBN)
	if err != nil {
		return fmt.Errorf("error deleting book from database: %w", err)
	}

	return nil
}

func (d *MySQLBookDAO) Update(book *models.Book) error {
	query := "UPDATE Books SET State = ?, OnHoldCustomerID = ?, CheckedOutCustomerID = ?, TimeUpdated = ? WHERE ISBN = ?"

	_, err := d.db.Exec(query, book.State, book.OnHoldCustomerID, book.CheckedOutCustomerID, book.TimeUpdated, book.ISBN)
	if err != nil {
		return fmt.Errorf("error updating book: %w", err)
	}

	return nil
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
			return nil, nil
		}
		
		return nil, fmt.Errorf("error: %w", err)
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
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	retrievedBooks := make([]*models.Book, 0)

	retrievedISBN := new(sql.NullString)
	retrievedState := new(sql.NullString)
	retrievedOnHoldCustomerID := new(sql.NullString)
	retrievedCheckedOutCustomerID := new(sql.NullString)
	retrievedTimeCreated := new(sql.NullString)
	retrievedTimeUpdated := new(sql.NullString)

	for rows.Next() {
		err := rows.Scan(
			retrievedISBN,
			retrievedState,
			retrievedOnHoldCustomerID,
			retrievedCheckedOutCustomerID,
			retrievedTimeCreated,
			retrievedTimeUpdated,
		)
			
		if err != nil {
			return nil, fmt.Errorf("error: %w", err)
		}

		nextBook := &models.Book{
			ISBN: nil,
			State: nil,
			OnHoldCustomerID: nil,
			CheckedOutCustomerID: nil,
			TimeCreated: nil,
			TimeUpdated: nil,
		}

		if retrievedISBN.Valid {
			nextBook.ISBN = &retrievedISBN.String
		}
	
		if retrievedState.Valid {
			nextBook.State = &retrievedState.String
		}
	
		if retrievedOnHoldCustomerID.Valid {
			nextBook.OnHoldCustomerID = &retrievedOnHoldCustomerID.String
		}
	
		if retrievedCheckedOutCustomerID.Valid {
			nextBook.CheckedOutCustomerID = &retrievedCheckedOutCustomerID.String
		}
	
		if retrievedTimeCreated.Valid {
			timeCreated, err := time.Parse("2006-01-02 15:04:05", retrievedTimeCreated.String)
			if err != nil {
				return nil, fmt.Errorf("error parsing time created in read: %w", err)
			}
			nextBook.TimeCreated = &timeCreated
		}
	
		if retrievedTimeUpdated.Valid {
			timeUpdated, err := time.Parse("2006-01-02 15:04:05", retrievedTimeUpdated.String)
			if err != nil {
				return nil, fmt.Errorf("error parsing time created in read: %w", err)
			}
			nextBook.TimeUpdated = &timeUpdated
		}

		retrievedBooks = append(retrievedBooks, nextBook)
	}

	return retrievedBooks, nil
}