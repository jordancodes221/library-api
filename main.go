package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	// "fmt"
)

type book struct{
	ISBN 	string 	`json:"isbn"`
	State 	string 	`json:"state"` // available, checked-out
}

var books = []book{
	{ISBN: "0000", State: "available"},
	{ISBN: "0001", State: "checked-out"},
}

// helper function
func bookByISBN(isbn string) (*book, error) {
	for i, b := range books {
		if b.ISBN == isbn {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

// GET ALL
func getAllBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}


// GET INDIVIDUAL
func getIndividualBook(c *gin.Context) {
	isbn := c.Param("isbn")
	book, err := bookByISBN(isbn)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// POST
func postBook(c *gin.Context) {
	var newBook book
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// CHECKOUT BOOK
func checkoutBook(c *gin.Context) {
	isbn, ok := c.GetQuery("isbn")

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Query parameter (ISBN) is missing."})
		return
	}

	book, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	if book.State == "checked-out" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book is already checked-out."})
		return
	}

	book.State = "checked-out"
	c.IndentedJSON(http.StatusOK, book)
}

// RETURN BOOK
func returnBook(c *gin.Context) {
	isbn, ok := c.GetQuery("isbn")
	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Query parameter (ISBN) is missing."})
		return
	}

	book, err := bookByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	if book.State == "available" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book has already been returned."})
		return
	}

	book.State = "available"
	c.IndentedJSON(http.StatusOK, book)
}

// DELETE
func deleteBook(c* gin.Context) {
	isbn, ok := c.GetQuery("isbn")
	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Query parameter (ISBN) is missing."})
		return
	}

	for i, b := range books {
		if b.ISBN == isbn {
			books = append(books[:i], books[i+1:]...)
		}
	}

	c.IndentedJSON(http.StatusOK, books)
}

func main() {
	router := gin.Default()
	router.GET("/books", getAllBooks)
	router.GET("/books/:isbn", getIndividualBook)
	router.POST("/books", postBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.DELETE("/delete", deleteBook)

	router.Run("localhost:8080")
}