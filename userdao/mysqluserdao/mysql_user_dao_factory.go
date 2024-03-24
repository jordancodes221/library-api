package mysqluserdao

import "example/library_project/userdao"

// "database/sql"

type MySQLUserDAOFactory struct {
	// db *sql.DB
	// dbUsername string
	// dbPassword string
	// dbHost string
	// dbPort string
	// dbName string
}

func NewMySQLUserDAOFactory(dbUsername string, dbPassword string, dbHost string, dbPort string, dbName string) *MySQLUserDAOFactory

func (f *MySQLUserDAOFactory) Open() error

func (f *MySQLUserDAOFactory) Close() error 

func (f *MySQLUserDAOFactory) BookDAO() userdao.UserDAO

func (f *MySQLUserDAOFactory) Clear() error