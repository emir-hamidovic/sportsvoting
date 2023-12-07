package goatplayers

import (
	"fmt"
	"log"
	"regexp"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"sportsvoting/request"
	"sportsvoting/scraper"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetGoatPlayersList() map[string]bool {
	fmt.Println("Getting goat players list")
	playerIDs := make(map[string]bool)

	urls := []string{
		"https://www.basketball-reference.com/leaders/pts_per_g_career.html",
		"https://www.basketball-reference.com/leaders/per_career.html",
		"https://www.basketball-reference.com/leaders/orb_pct_career.html",
		"https://www.basketball-reference.com/leaders/dbpm_career.html",
		"https://www.basketball-reference.com/leaders/pts_per_g_career_p.html",
		"https://www.basketball-reference.com/leaders/per_career_p.html",
		"https://www.basketball-reference.com/leaders/orb_pct_career_p.html",
		"https://www.basketball-reference.com/leaders/bpm_career_p.html",
		"https://www.basketball-reference.com/leaders/def_rtg_career_p.html",
		"https://www.basketball-reference.com/leaders/trb_per_g_career_p.html",
	}

	for _, url := range urls {
		doc, err := request.GetDocumentFromURL(url)
		if err != nil {
			log.Println(err)
		}

		doc.Find("table#nba > tbody > tr").Each(func(i int, row *goquery.Selection) {
			playerID := request.GetPlayerIDFromLeadersDocument(row)
			if playerID != "" && !playerIDs[playerID] {
				playerIDs[playerID] = true
			}
		})

		time.Sleep(3 * time.Second)
	}

	return playerIDs
}

func InsertGoatPlayerStats(playerIds map[string]bool, db database.Database) {
	for playerID := range playerIds {
		goatPlayer, goatStatsRegular, goatStatsPlayoffs := scrapePlayerInfo(playerID)

		if goatPlayer.ID != "" {
			_, err := db.InsertGOATPlayer(goatPlayer)
			if err != nil {
				fmt.Println(err)
				time.Sleep(4 * time.Second)
				continue
			}

			_, err = db.InsertGOATStats(goatStatsRegular)
			if err != nil {
				fmt.Println(err)
				time.Sleep(4 * time.Second)
				continue
			}

			_, err = db.InsertGOATStats(goatStatsPlayoffs)
			if err != nil {
				fmt.Println(err)
				time.Sleep(4 * time.Second)
				continue
			}

			time.Sleep(4 * time.Second)
		}
	}
}

func scrapePlayerInfo(playerID string) (databasestructs.GoatPlayers, databasestructs.GoatStats, databasestructs.GoatStats) {
	var goatPlayer databasestructs.GoatPlayers
	var goatStatsRegular, goatStatsPlayoffs databasestructs.GoatStats
	url := fmt.Sprintf("https://www.basketball-reference.com/players/%s/%s.html", string(playerID[0]), playerID)
	fmt.Println(url)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		log.Println(err)
		time.Sleep(4 * time.Second)
		return goatPlayer, goatStatsRegular, goatStatsPlayoffs
	}

	goatPlayer.ID = playerID
	goatStatsRegular.PlayerID = playerID
	goatStatsPlayoffs.PlayerID = playerID
	goatStatsPlayoffs.IsPlayoffs = true
	doc.Find("div#meta h1 span").Each(func(i int, s *goquery.Selection) {
		goatPlayer.Name = s.Text()
	})
	fmt.Println(goatPlayer.Name)

	doc.Find("div#meta p").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Experience") {
			goatPlayer.IsActive = true
			goatStatsPlayoffs.IsActive = true
			goatStatsRegular.IsActive = true
		}
	})

	goatStatsPlayoffs.Position = scraper.GetTDDataStatString(doc.Find("table#per_game tbody tr:first-child"), "pos")
	goatStatsRegular.Position = goatStatsPlayoffs.Position

	doc.Find("table#per_game tfoot tr").Each(func(i int, s *goquery.Selection) {
		league := s.Find("td[data-stat='lg_id']").Text()
		if league == "NBA" && goatStatsRegular.Points == 0 {
			goatStatsRegular.Points = scraper.GetTDDataStatFloat(s, "pts_per_g")
			goatStatsRegular.Rebounds = scraper.GetTDDataStatFloat(s, "trb_per_g")
			goatStatsRegular.Assists = scraper.GetTDDataStatFloat(s, "ast_per_g")
			goatStatsRegular.Blocks = scraper.GetTDDataStatFloat(s, "blk_per_g")
			goatStatsRegular.Steals = scraper.GetTDDataStatFloat(s, "stl_per_g")
			goatStatsRegular.Turnovers = scraper.GetTDDataStatFloat(s, "tov_per_g")
			goatStatsRegular.FGPercentage = scraper.GetTDDataStatFloat(s, "fg_pct") * 100
			goatStatsRegular.ThreeFGPercentage = scraper.GetTDDataStatFloat(s, "fg3_pct") * 100
			goatStatsRegular.FTPercentage = scraper.GetTDDataStatFloat(s, "ft_pct") * 100
		}
	})

	doc.Find("table#playoffs_per_game tfoot tr").Each(func(i int, s *goquery.Selection) {
		league := s.Find("td[data-stat='lg_id']").Text()
		if league == "NBA" && goatStatsPlayoffs.Points == 0 {
			goatStatsPlayoffs.Points = scraper.GetTDDataStatFloat(s, "pts_per_g")
			goatStatsPlayoffs.Rebounds = scraper.GetTDDataStatFloat(s, "trb_per_g")
			goatStatsPlayoffs.Assists = scraper.GetTDDataStatFloat(s, "ast_per_g")
			goatStatsPlayoffs.Blocks = scraper.GetTDDataStatFloat(s, "blk_per_g")
			goatStatsPlayoffs.Steals = scraper.GetTDDataStatFloat(s, "stl_per_g")
			goatStatsPlayoffs.Turnovers = scraper.GetTDDataStatFloat(s, "tov_per_g")
			goatStatsPlayoffs.FGPercentage = scraper.GetTDDataStatFloat(s, "fg_pct") * 100
			goatStatsPlayoffs.ThreeFGPercentage = scraper.GetTDDataStatFloat(s, "fg3_pct") * 100
			goatStatsPlayoffs.FTPercentage = scraper.GetTDDataStatFloat(s, "ft_pct") * 100
		}
	})

	doc.Find("table#totals tfoot tr").Each(func(i int, s *goquery.Selection) {
		league := s.Find("td[data-stat='lg_id']").Text()
		if league == "NBA" && goatStatsRegular.TotalPoints == 0 {
			goatStatsRegular.TotalPoints = scraper.GetTDDataStatInt(s, "pts")
			goatStatsRegular.TotalRebounds = scraper.GetTDDataStatInt(s, "trb")
			goatStatsRegular.TotalAssists = scraper.GetTDDataStatInt(s, "ast")
			goatStatsRegular.TotalSteals = scraper.GetTDDataStatInt(s, "stl")
			goatStatsRegular.TotalBlocks = scraper.GetTDDataStatInt(s, "blk")
		}
	})

	doc.Find("table#playoffs_totals tfoot tr").Each(func(i int, s *goquery.Selection) {
		league := s.Find("td[data-stat='lg_id']").Text()
		if league == "NBA" && goatStatsPlayoffs.TotalPoints == 0 {
			goatStatsPlayoffs.TotalPoints = scraper.GetTDDataStatInt(s, "pts")
			goatStatsPlayoffs.TotalRebounds = scraper.GetTDDataStatInt(s, "trb")
			goatStatsPlayoffs.TotalAssists = scraper.GetTDDataStatInt(s, "ast")
			goatStatsPlayoffs.TotalSteals = scraper.GetTDDataStatInt(s, "stl")
			goatStatsPlayoffs.TotalBlocks = scraper.GetTDDataStatInt(s, "blk")
		}
	})

	doc.Find("table#advanced tfoot tr").Each(func(i int, s *goquery.Selection) {
		league := s.Find("td[data-stat='lg_id']").Text()
		if league == "NBA" && goatStatsRegular.PER == 0 {
			goatStatsRegular.PER = scraper.GetTDDataStatFloat(s, "per")
			goatStatsRegular.OffBPM = scraper.GetTDDataStatFloat(s, "obpm")
			goatStatsRegular.DefBPM = scraper.GetTDDataStatFloat(s, "dbpm")
			goatStatsRegular.BPM = scraper.GetTDDataStatFloat(s, "bpm")
			goatStatsRegular.VORP = scraper.GetTDDataStatFloat(s, "vorp")
			goatStatsRegular.OffWS = scraper.GetTDDataStatFloat(s, "ows")
			goatStatsRegular.DefWS = scraper.GetTDDataStatFloat(s, "dws")
			goatStatsRegular.WS = scraper.GetTDDataStatFloat(s, "ws")
		}
	})

	doc.Find("table#playoffs_advanced tfoot tr").Each(func(i int, s *goquery.Selection) {
		league := s.Find("td[data-stat='lg_id']").Text()
		if league == "NBA" && goatStatsPlayoffs.PER == 0 {
			goatStatsPlayoffs.PER = scraper.GetTDDataStatFloat(s, "per")
			goatStatsPlayoffs.OffBPM = scraper.GetTDDataStatFloat(s, "obpm")
			goatStatsPlayoffs.DefBPM = scraper.GetTDDataStatFloat(s, "dbpm")
			goatStatsPlayoffs.BPM = scraper.GetTDDataStatFloat(s, "bpm")
			goatStatsPlayoffs.VORP = scraper.GetTDDataStatFloat(s, "vorp")
			goatStatsPlayoffs.OffWS = scraper.GetTDDataStatFloat(s, "ows")
			goatStatsPlayoffs.DefWS = scraper.GetTDDataStatFloat(s, "dws")
			goatStatsPlayoffs.WS = scraper.GetTDDataStatFloat(s, "ws")
		}
	})

	doc.Find("table#per_poss tfoot tr").Each(func(i int, s *goquery.Selection) {
		league := s.Find("td[data-stat='lg_id']").Text()
		if league == "NBA" && goatStatsRegular.DefRtg == 0 {
			goatStatsRegular.DefRtg = scraper.GetTDDataStatFloat(s, "def_rtg")
			goatStatsRegular.OffRtg = scraper.GetTDDataStatFloat(s, "off_rtg")
		}
	})

	doc.Find("table#playoffs_per_poss tfoot tr").Each(func(i int, s *goquery.Selection) {
		league := s.Find("td[data-stat='lg_id']").Text()
		if league == "NBA" && goatStatsPlayoffs.DefRtg == 0 {
			goatStatsPlayoffs.DefRtg = scraper.GetTDDataStatFloat(s, "def_rtg")
			goatStatsPlayoffs.OffRtg = scraper.GetTDDataStatFloat(s, "off_rtg")
		}
	})

	doc.Find("ul#bling li").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Finals MVP") {
			goatPlayer.FMVP = parseAccolade(s, "Finals MVP")
		} else if strings.Contains(s.Text(), "MVP") {
			goatPlayer.MVP = parseAccolade(s, "MVP")
		} else if strings.Contains(s.Text(), "ROY") {
			goatPlayer.ROY = parseAccolade(s, "ROY")
		} else if strings.Contains(s.Text(), "Def. POY") {
			goatPlayer.Dpoy = parseAccolade(s, "Def. POY")
		} else if strings.Contains(s.Text(), "All-NBA") {
			goatPlayer.AllNba = parseAccolade(s, "All-NBA")
		} else if strings.Contains(s.Text(), "All Star") {
			goatPlayer.AllStar = parseAccolade(s, "All Star")
		} else if strings.Contains(s.Text(), "NBA Champ") {
			goatPlayer.Championships = parseAccolade(s, "NBA Champ")
		} else if strings.Contains(s.Text(), "All-Defensive") {
			goatPlayer.AllDefense = parseAccolade(s, "All-Defensive")
		}
	})

	return goatPlayer, goatStatsRegular, goatStatsPlayoffs
}

func UpdateActiveGOATStats(db database.Database) error {
	rows, err := db.GetActivePlayers()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var playerID string
		err := rows.Scan(&playerID)
		if err != nil {
			return err
		}

		goatPlayer, goatStatsRegular, goatStatsPlayoffs := scrapePlayerInfo(playerID)
		if goatPlayer.ID != "" {
			_, err = db.UpdateGOATPlayer(goatPlayer)
			if err != nil {
				fmt.Println(err)
				time.Sleep(4 * time.Second)
				continue
			}

			_, err := db.UpdateGOATStats(goatStatsRegular)
			if err != nil {
				fmt.Println(err)
			}

			_, err = db.UpdateGOATStats(goatStatsPlayoffs)
			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(4 * time.Second)
		}
	}

	return nil
}

func parseAccolade(s *goquery.Selection, key string) int64 {
	re := regexp.MustCompile(fmt.Sprintf(`(\d+)x %s`, key))
	match := re.FindStringSubmatch(s.Text())
	if len(match) > 1 {
		res, _ := strconv.ParseInt(match[1], 10, 64)
		return res
	} else {
		return 1
	}
}
