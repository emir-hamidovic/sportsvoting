package advancedstats

import (
	"fmt"
	"sportsvoting/database"
	"sportsvoting/request"
	"sportsvoting/scraper"

	"github.com/PuerkitoBio/goquery"
)

type AdvancedStats struct {
	PlayerID string  `json:"stats,omitempty"`
	TeamAbbr string  `json:"team,omitempty"`
	Season   string  `json:"season,omitempty"`
	PER      float64 `json:"per,omitempty"`
	TSPct    float64 `json:"ts,omitempty"`
	USGPCt   float64 `json:"usg,omitempty"`
	OffWS    float64 `json:"ows,omitempty"`
	DefWS    float64 `json:"dws,omitempty"`
	WS       float64 `json:"ws,omitempty"`
	OffBPM   float64 `json:"obpm,omitempty"`
	DefBPM   float64 `json:"dbpm,omitempty"`
	BPM      float64 `json:"bpm,omitempty"`
	VORP     float64 `json:"vorp,omitempty"`
	DefRtg   float64 `json:"defrtg,omitempty"`
	OffRtg   float64 `json:"offrtg,omitempty"`
}

func FillPlayerStatsForSeason(row *goquery.Selection, season string, stats *AdvancedStats) {
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

func UpdateStats(db database.Database, stats AdvancedStats) error {
	// Problem with update, check this
	fmt.Println(stats)
	res, err := db.UpdateAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.OffRtg, stats.DefRtg, stats.TeamAbbr, stats.PlayerID, stats.Season)
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
		_, err = db.InsertAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.OffRtg, stats.DefRtg, stats.TeamAbbr, stats.PlayerID, stats.Season)
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

			var stats AdvancedStats
			FillPlayerStatsForSeason(row, season, &stats)
			_, err := db.UpdateTradedPlayerAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.Season, id)
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
