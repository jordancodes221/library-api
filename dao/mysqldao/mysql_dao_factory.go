package mysqldao

import (
	"example/library_project/dao"
	// "example/library_project/models"

	_ "github.com/go-sql-driver/mysql"
	"database/sql"

	"fmt"
	// "log"

	// "os"
)

type MySQLDAOFactory struct {
	db *sql.DB
	dbUsername string
	dbPassword string
	dbHost string
	dbPort string
	dbName string
}

func NewMySQLDAOFactory(dbUsername string, dbPassword string, dbHost string, dbPort string, dbName string) *MySQLDAOFactory {
	return &MySQLDAOFactory{
		db: nil,
		dbUsername: dbUsername,
		dbPassword: dbPassword,
		dbHost: dbHost,
		dbPort: dbPort,
		dbName: dbName,
	}
}

func (f *MySQLDAOFactory) Open() error {
	dbUsername := f.dbUsername
	dbPassword := f.dbPassword
	dbHost := f.dbHost
	dbPort := f.dbPort
	dbName := f.dbName

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping the database: %w", err)
	}

	f.db = db

	// log.Println("Connected to the MySQL database")
	fmt.Println("Connected to the MySQL database")

	return nil
}

func (f *MySQLDAOFactory) Close() error {
	if f.db == nil {
		return nil
	}

	err := f.db.Close()
	if err != nil {
		return fmt.Errorf("Failed to close database connection: %w", err)
	}

	return nil
}

func (f *MySQLDAOFactory) BookDAO() dao.BookDAO {
	return &MySQLBookDAO{
		db: f.db,
	}
}

func (f *MySQLDAOFactory) Clear() error {
	_, err := f.db.Exec("TRUNCATE TABLE Books;")
	if err != nil {
		return fmt.Errorf("failed to clear database: %w", err)
	}

	return nil
}