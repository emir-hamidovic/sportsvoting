package advancedstats

import (
	"fmt"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"sportsvoting/request"
	"sportsvoting/scraper"

	"github.com/PuerkitoBio/goquery"
)

func FillPlayerStatsForSeason(row *goquery.Selection, season string, stats *databasestructs.AdvancedStats) {
	stats.PER = scraper.GetTDDataStatFloat(row, "per")
	stats.TSPct = scraper.GetTDDataStatFloat(row, "ts_pct")
	stats.USGPCt = scraper.GetTDDataStatFloat(row, "usg_pct")
	stats.OffWS = scraper.GetTDDataStatFloat(row, "ows")
	stats.DefWS = scraper.GetTDDataStatFloat(row, "dws")
	stats.WS = scraper.GetTDDataStatFloat(row, "ws")
	stats.OffBPM = scraper.GetTDDataStatFloat(row, "obpm")
	stats.DefBPM = scraper.GetTDDataStatFloat(row, "dbpm")
	stats.BPM = scraper.GetTDDataStatFloat(row, "bpm")
	stats.VORP = scraper.GetTDDataStatFloat(row, "vorp")
	stats.Season = season
}

func UpdateStats(db database.Database, stats databasestructs.AdvancedStats) error {
	// Problem with update, check this
	fmt.Println(stats)
	res, err := db.UpdateAdvancedStats(stats)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rows, stats.Season)

	if rows == 0 && stats.Season != "" {
		fmt.Println("insert advanced", stats)
		_, err = db.InsertAdvancedStats(stats)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Player advanced stats added to database.")
	return nil
}

func UpdateTradedPlayerStats(db database.Database, season string) error {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_advanced.html", season)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows := doc.Find("table#advanced_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		team := scraper.GetTDDataStatString(row, "team_id")
		if team == "TOT" {
			id := request.GetPlayerIDFromDocument(row)

			var stats databasestructs.AdvancedStats
			FillPlayerStatsForSeason(row, season, &stats)
			stats.PlayerID = id
			_, err := db.UpdateTradedPlayerAdvancedStats(stats)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	url = fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_poss.html", season)
	doc, err = request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows = doc.Find("table#per_poss_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		team := scraper.GetTDDataStatString(row, "team_id")
		if team == "TOT" {
			id := request.GetPlayerIDFromDocument(row)

			defrtg := scraper.GetTDDataStatFloat(row, "def_rtg")
			offrtg := scraper.GetTDDataStatFloat(row, "off_rtg")
			_, err := db.UpdateOffAndDefRtg(offrtg, defrtg, season, id)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	return nil
}
