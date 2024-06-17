package mysqluser

import "github.com/mohsenHa/messenger/repository/mysql"

type DB struct {
	conn *mysql.DB
}

func New(conn *mysql.DB) *DB {
	return &DB{
		conn: conn,
	}
}
