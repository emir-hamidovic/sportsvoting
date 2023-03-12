package main

import (
	"context"
	"errors"
	"log"
	"scraper/database"
	"scraper/parser/players"
	"scraper/parser/teams"
	"sync"
	"time"
)

func Parse(db database.Database, parse string) error {
	if parse == "teams" {
		return teams.ParseTeams(db)
	} else if parse == "players" {
		return players.ParsePlayers(db)
	}

	return errors.New("incorrect type")
}

func RunUpdate(db database.Database, ctx context.Context, errCh chan<- error) {
	ticker := time.NewTicker(24 * time.Hour)
	var err error
	for {
		select {
		case <-ctx.Done():
			return // if the context is cancelled, return without sending the error
		case <-ticker.C:
			now := time.Now().UTC()
			if now.Hour() == 8 && now.Minute() == 0 {
				err = players.ParseBoxScores(db)
				if err != nil {
					errCh <- err
					return
				}
			}
		}
	}
}

func main() {
	db, err := database.NewDB(database.Config{DbType: "mysql", DbName: "nba", Addr: "localhost:3306"})
	if err != nil {
		log.Fatal(err)
	}

	err = players.ParseBoxScores(db)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Use sync.Once to execute a function only once.
	var once sync.Once

	errChan := make(chan error, 1)

	go func(ctx context.Context, errCh chan<- error) {
		once.Do(func() {
			err = Parse(db, "teams")
			if err != nil {
				return
			}

			err = Parse(db, "players")
			if err != nil {
				return
			}
		})

		select {
		case <-ctx.Done():
			return // if the context is cancelled, return without sending the error
		case errCh <- err:
			return
		}
	}(ctx, errChan)

	go RunUpdate(db, ctx, errChan)

	go func(errCh <-chan error) {
		for err := range errCh {
			if err != nil {
				log.Println(err)
				cancel() // cancel the context if an error is received
			}
		}
	}(errChan)

	<-ctx.Done()
	// Close the error channel.
	close(errChan)
}
