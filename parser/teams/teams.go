package teams

import (
	"fmt"
	"scraper/database"
	"scraper/parser"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Teams struct {
	TeamAbbr         string
	Name             string
	Logo             string
	WinLossPct       float64
	Playoffs         int64
	DivisionTitles   int64
	ConferenceTitles int64
	Championships    int64
}

func ParseTeams(db database.Database) error {
	allTeams, err := ParseAll()
	if err != nil {
		return err
	}

	stmt, err := db.PrepareStatementForTeamInsert()
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, team := range allTeams {
		_, err = stmt.Exec(team.TeamAbbr, team.Name, team.Logo, team.WinLossPct, team.Playoffs, team.DivisionTitles, team.ConferenceTitles, team.Championships)
		if err != nil {
			return err
		}
	}

	fmt.Println("Teams added to database.")
	return nil
}

func ParseAll() ([]Teams, error) {
	teams := []Teams{}
	url := "https://www.basketball-reference.com/teams"
	res, err := parser.SendRequest(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	teams = findBasicTeamInfo(doc, teams)

	// Wait for 5 seconds before sending the next request to try avoiding the rate limit of 20req/min
	time.Sleep(5 * time.Second)

	return teams, nil
}

func findBasicTeamInfo(doc *goquery.Document, teams []Teams) []Teams {
	rows := doc.Find("table#teams_active > tbody > tr.full_table")
	rows.Each(func(i int, row *goquery.Selection) {
		name := row.Find("th[data-stat='franch_name']").Text()
		abbr, exists := row.Find("th[data-stat='franch_name'] > a").Attr("href")
		if exists {
			// Extract team ID from URL
			idParts := strings.Split(abbr, "/")
			if len(idParts) > 2 {
				abbr = strings.TrimSuffix(idParts[2], ".html")
			}

			winlosspct, _ := strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='win_loss_pct']").Text()), 64)
			playoffs, _ := strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='years_playoffs']").Text()), 10, 64)
			divtitles, _ := strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='years_division_champion']").Text()), 10, 64)
			conftitles, _ := strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='years_conference_champion']").Text()), 10, 64)
			championships, _ := strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='years_league_champion']").Text()), 10, 64)
			logo, err := getLogoForTeam(abbr)
			if err != nil {
				return
			}

			teams = append(teams, Teams{Name: name, TeamAbbr: abbr, WinLossPct: winlosspct * 100, Playoffs: playoffs, DivisionTitles: divtitles, ConferenceTitles: conftitles, Championships: championships, Logo: logo})
			time.Sleep(4 * time.Second)
		}
	})

	return teams
}

func getLogoForTeam(abbr string) (string, error) {
	url := fmt.Sprintf("https://www.basketball-reference.com/teams/%s", abbr)
	res, err := parser.SendRequest(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	logo, _ := doc.Find("img.teamlogo").Attr("src")

	return logo, nil
}
