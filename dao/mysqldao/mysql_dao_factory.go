package mysqldao

import (
	"example/library_project/dao"
	// "example/library_project/models"

	_ "github.com/go-sql-driver/mysql"
	"database/sql"

	"fmt"
	// "log"
)

type MySQLDAOFactory struct {
	db *sql.DB
}

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/Library")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
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