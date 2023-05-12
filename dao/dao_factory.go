package dao

type DAOFactory interface{
	BookDAO() BookDAO
}