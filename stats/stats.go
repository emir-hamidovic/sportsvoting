package stats

import (
	"fmt"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"sportsvoting/request"
	"sportsvoting/scraper"

	"github.com/PuerkitoBio/goquery"
)

func FillPlayerStatsForSeason(row *goquery.Selection, season string, stats *databasestructs.PlayerStats) {
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

func UpdateStats(db database.Database, stats databasestructs.PlayerStats) error {
	fmt.Println(stats)
	res, err := db.UpdateStats(stats)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	if rows == 0 && stats.Season != "" {
		_, err = db.InsertStats(stats)
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

			var stats databasestructs.PlayerStats
			FillPlayerStatsForSeason(row, season, &stats)
			stats.Position = scraper.GetTDDataStatString(row, "pos")
			stats.PlayerID = id
			_, err := db.UpdateTradedPlayerStats(stats)
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
