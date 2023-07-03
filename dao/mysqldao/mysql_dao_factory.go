package mysqldao

import (
	"example/library_project/dao"
	// "example/library_project/models"

	_ "github.com/go-sql-driver/mysql"
	"database/sql"

	"fmt"
	// "log"

	"os"
)

type MySQLDAOFactory struct {
	db *sql.DB
}

func ConnectDB() (*sql.DB, error) {
	dbUsername, ok := os.LookupEnv("LIBRARY_DB_USERNAME")
	if !ok {
		return nil, fmt.Errorf("Error retrieving LIBRARY_DB_USERNAME environment variable.")
	}

	dbPassword, ok := os.LookupEnv("LIBRARY_DB_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("Error retrieving LIBRARY_DB_PASSWORD environment variable.")
	}

	dbHost, ok := os.LookupEnv("LIBRARY_DB_HOST")
	if !ok {
		return nil, fmt.Errorf("Error retrieving LIBRARY_DB_HOST environment variable.")
	}

	dbPort, ok := os.LookupEnv("LIBRARY_DB_PORT")
	if !ok {
		return nil, fmt.Errorf("Error retrieving LIBRARY_DB_PORT environment variable.")
	}

	dbName, ok := os.LookupEnv("LIBRARY_DB_NAME")
	if !ok {
		return nil, fmt.Errorf("Error retrieving LIBRARY_DB_NAME environment variable.")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	// log.Println("Connected to the MySQL database")
	fmt.Println("Connected to the MySQL database")

	return db, nil
}

func NewMySQLDAOFactory() (*MySQLDAOFactory, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	factory := MySQLDAOFactory{
		db: db,
	}

	return &factory, nil
}

func (f *MySQLDAOFactory) BookDAO() dao.BookDAO {
	return &MySQLBookDAO{
		db: f.db,
	}
}