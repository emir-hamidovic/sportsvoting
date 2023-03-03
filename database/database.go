package database

import (
	"database/sql"
	"errors"
	"scraper/database/mysql_db"
)

const (
	MYSQL string = "mysql"
)

type Database interface {
	InsertPlayers()
	InsertTeams()
	PrepareStatementForPlayerInsert() (*sql.Stmt, error)
	PrepareStatementForTeamInsert() (*sql.Stmt, error)
}

type Config struct {
	DbType string
	DbName string
	Addr   string
}

func NewDB(conf Config) (Database, error) {
	switch conf.DbType {
	case MYSQL:
		return mysql_db.NewDB(conf.DbName, conf.Addr)
	}

	return nil, errors.New("incorrect db type entered")
}
