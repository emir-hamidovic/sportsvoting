package stats

import (
	"fmt"
	"sportsvoting/database"
	"sportsvoting/request"
	"sportsvoting/scraper"

	"github.com/PuerkitoBio/goquery"
)

type Stats struct {
	PlayerID          string  `json:"playerid,omitempty"`
	Games             int64   `json:"g,omitempty"`
	GamesStarted      int64   `json:"gs,omitempty"`
	Minutes           float64 `json:"mpg,omitempty"`
	Points            float64 `json:"ppg,omitempty"`
	Rebounds          float64 `json:"rpg,omitempty"`
	Assists           float64 `json:"apg,omitempty"`
	Steals            float64 `json:"spg,omitempty"`
	Blocks            float64 `json:"bpg,omitempty"`
	Turnovers         float64 `json:"topg,omitempty"`
	FGPercentage      float64 `json:"fgpct,omitempty"`
	ThreeFGPercentage float64 `json:"threefgpct,omitempty"`
	FTPercentage      float64 `json:"ftpct,omitempty"`
	Season            string  `json:"season,omitempty"`
	Position          string  `json:"position,omitempty"`
	TeamAbbr          string  `json:"team,omitempty"`
	IsRookie          bool    `json:"rookie,omitempty"`
}

func FillPlayerStatsForSeason(row *goquery.Selection, season string, stats *Stats) {
	stats.Games = scraper.GetTDDataStatInt(row, "g")
	stats.GamesStarted = scraper.GetTDDataStatInt(row, "gs")
	stats.Minutes = scraper.GetTDDataStatFloat(row, "mp_per_g")
	stats.Points = scraper.GetTDDataStatFloat(row, "pts_per_g")
	stats.Rebounds = scraper.GetTDDataStatFloat(row, "trb_per_g")
	stats.Assists = scraper.GetTDDataStatFloat(row, "ast_per_g")
	stats.Blocks = scraper.GetTDDataStatFloat(row, "blk_per_g")
	stats.Steals = scraper.GetTDDataStatFloat(row, "stl_per_g")
	stats.Turnovers = scraper.GetTDDataStatFloat(row, "tov_per_g")
	stats.FGPercentage = scraper.GetTDDataStatFloat(row, "fg_pct") * 100
	stats.ThreeFGPercentage = scraper.GetTDDataStatFloat(row, "fg3_pct") * 100
	stats.FTPercentage = scraper.GetTDDataStatFloat(row, "ft_pct") * 100
	stats.Season = season
}

func UpdateStats(db database.Database, stats Stats) error {
	fmt.Println(stats)
	res, err := db.UpdateStats(stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID, stats.TeamAbbr)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	if rows == 0 && stats.Season != "" {
		_, err = db.InsertStats(stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID, stats.TeamAbbr)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Player stats added to database.")
	return nil
}

func UpdateTradedPlayerStats(db database.Database, season string) error {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_game.html", season)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows := doc.Find("table#per_game_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		team := scraper.GetTDDataStatString(row, "team_id")
		if team == "TOT" {
			id := request.GetPlayerIDFromDocument(row)

			var stats Stats
			FillPlayerStatsForSeason(row, season, &stats)
			position := scraper.GetTDDataStatString(row, "pos")
			_, err := db.UpdateTradedPlayerStats(stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, position, id)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	return nil
}

func SetRookies(db database.Database, season string) error {
	fmt.Println("Setting rookies")
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_rookies.html", season)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows := doc.Find("table#rookies > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		id := request.GetPlayerIDFromDocument(row)
		_, err := db.SetRookieStatus(id)
		if err != nil {
			fmt.Println(err)
		}
	})

	return nil
}
