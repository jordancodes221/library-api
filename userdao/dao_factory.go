package userdao

type DAOFactory interface{
	UserDAO() UserDAO
	Open() error
	Close() error
	Clear() error
}