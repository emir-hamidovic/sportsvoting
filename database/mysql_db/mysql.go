package mysql_db

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

type MySqlDB struct {
	db *sql.DB
}

func NewDB(dbname string, addr string) (*MySqlDB, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               dbname,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return &MySqlDB{db: db}, nil
}
