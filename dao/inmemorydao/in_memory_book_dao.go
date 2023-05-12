package inmemorydao

import (
	"example/library_project/models"
)

type InMemoryBookDAO struct {
	Books map[string]*models.Book
}

func (d *InMemoryBookDAO) Create(newBook *models.Book) {
	d.Books[*newBook.ISBN] = newBook
}

func (d *InMemoryBookDAO) Delete(book *models.Book) {
	delete(d.Books, *book.ISBN)
}

func (d *InMemoryBookDAO) Update(book *models.Book) {
	d.Books[*book.ISBN] = book
}

func (d *InMemoryBookDAO) Read(isbn string) (*models.Book, error) {
	retrievedBook, ok := d.Books[isbn] // in the future, this could be a call to a database

	// For scalability, we can add a database connection here. 
	// If there is an error connecting to the database, then we will return: nil, InternalServerError

	if ok {
		return retrievedBook, nil
	} else {
		return nil, nil
	}
}

func (d *InMemoryBookDAO) ReadAll() ([]*models.Book, error) {	
	all_books := make([]*models.Book, 0)

	// For scalability, we can add a database connection here. 
	// If there is an error connecting to the database, then we will return: nil, InternalServerError

	for _, currentBook := range d.Books {
		all_books = append(all_books, currentBook)
	}

	return all_books, nil
}