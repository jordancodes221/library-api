package dao

type DAOFactory interface{
	BookDAO() BookDAO
	Open() error
	Close() error
	Clear() error
}