package main

import (
	"example/library_project/handlers"
	// "example/library_project/models"
	"example/library_project/dao"
	"example/library_project/utils"

	"example/library_project/dao/inmemorydao"
	"example/library_project/dao/mysqldao"

	"example/library_project/testdata"

	// "net/http"
	"github.com/gin-gonic/gin"
	// "errors"
	// "time"
	// "encoding/json"
	"fmt"

	// "reflect"
	// "strconv"

	"log"

	"os"
	"os/signal"
	"syscall"
)

func main() {

	// DAO selection
	var daoFactory dao.DAOFactory
	daoSelection := os.Getenv("DAO_SELECTION")
	testMode := os.Getenv("TEST_MODE")

	if daoSelection == "inmemory" {
		daoFactory = inmemorydao.NewInMemoryDAOFactory()
	} else if daoSelection == "mysql" {
		// Environment variables for database
		dbUsername := os.Getenv("LIBRARY_DB_USERNAME")
		dbPassword := os.Getenv("LIBRARY_DB_PASSWORD")
		dbHost := os.Getenv("LIBRARY_DB_HOST")
		dbPort := os.Getenv("LIBRARY_DB_PORT")
		dbName := os.Getenv("LIBRARY_DB_NAME")

		daoFactory = mysqldao.NewMySQLDAOFactory(dbUsername, dbPassword, dbHost, dbPort, dbName)
	} else {
		log.Fatal("unexpected dao selection")
	}

	// Open the database connection
	if err := daoFactory.Open(); err != nil {
		log.Fatal("failed to open database connection: ", err)
	}

	// Defer closing the database connection
	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel, syscall.SIGINT)

	go func() {
		receivedSignal := <- signalsChannel
		fmt.Println("RECEIVED SIGNAL...")

		if receivedSignal == syscall.SIGINT {
			fmt.Println("SIGNAL IS SIGINT...")
			daoFactory.Close()
			fmt.Println("ABOUT TO EXIT...")
			os.Exit(0)
		}
	}()

	// Instantiate DAO
	bookDAO := daoFactory.BookDAO()

	// If in integration test mode, instantiate test data and add to database
	if testMode == "integration" {
		testBooks, err := testdata.InstantiateIntegrationTestData()
		if err != nil{
			log.Fatal("failed to instantiate test data")
		}

		for _, currentTestBook := range testBooks {
			if err := bookDAO.Create(currentTestBook); err != nil{
				log.Fatal("failed to add test data to DAO")
			}
		}
	}

	realTimeProvider := &utils.ProductionDateTimeProvider{}
	h := handlers.NewBooksHandler(bookDAO, realTimeProvider)

	router := gin.Default()
	router.GET("/books", h.GetAllBooks)
	router.GET("/books/:isbn", h.GetIndividualBook)
	router.POST("/books", h.CreateBook)
	router.DELETE("/books/:isbn", h.DeleteBook)
	router.PATCH("/books/:isbn", h.UpdateBook)

	fmt.Println("ABOUT TO CALL ROUTER.RUN...")
	router.Run("localhost:8080")
}