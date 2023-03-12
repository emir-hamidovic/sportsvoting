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
	InsertPlayer(string, string, string, string, string, string, int64, int64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64) (sql.Result, error)
	InsertTeam(string, string, string, float64, int64, int64, int64, int64) (sql.Result, error)
	UpdatePlayers(int64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, string) (sql.Result, error)
	SelectPlayerID(string) (*sql.Rows, error)
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
