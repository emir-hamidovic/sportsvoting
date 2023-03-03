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

func (m *MySqlDB) InsertPlayers() {}
func (m *MySqlDB) InsertTeams()   {}
func (m *MySqlDB) PrepareStatementForPlayerInsert() (*sql.Stmt, error) {
	return m.db.Prepare("INSERT INTO players(playerid, name, college, teamabbr, height, weight, age, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, turnoverspergame, fgpercentage, ftpercentage, threeptpercentage) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
}

func (m *MySqlDB) PrepareStatementForTeamInsert() (*sql.Stmt, error) {
	return m.db.Prepare("INSERT INTO teams(teamabbr, name, logo, winlosspct, playoffs, divisiontitles, conferencetitles, championships) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
}
