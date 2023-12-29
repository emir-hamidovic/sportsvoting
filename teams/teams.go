package teams

import (
	"fmt"
	"log"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"sportsvoting/players"
	"sportsvoting/request"
	"sportsvoting/scraper"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ParseTeams(db database.Database, season string) (map[string]databasestructs.PlayerInfo, error) {
	allTeams, err := ParseBasicInfoForEveryTeam()
	if err != nil {
		return nil, err
	}

	roster := make(map[string]databasestructs.PlayerInfo, 600)
	for _, team := range allTeams {
		url := fmt.Sprintf("https://www.basketball-reference.com/teams/%s/%s.html", team.TeamAbbr, season)
		doc, err := request.GetDocumentFromURL(url)
		if err != nil {
			return nil, err
		}

		roster, err = players.GetPlayerInfo(doc, team.TeamAbbr, roster, season)
		if err != nil {
			return nil, err
		}

		team.Logo = scraper.GetTeamLogo(doc)
		log.Println(team)

		_, err = db.InsertTeam(team)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(4 * time.Second)
	}

	log.Println("Teams added to database.")
	return roster, nil
}

func ParseBasicInfoForEveryTeam() ([]databasestructs.TeamInfo, error) {
	url := "https://www.basketball-reference.com/teams"
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return nil, err
	}

	return findBasicTeamInfo(doc, []databasestructs.TeamInfo{}), nil
}

func findBasicTeamInfo(doc *goquery.Document, teams []databasestructs.TeamInfo) []databasestructs.TeamInfo {
	rows := doc.Find("table#teams_active > tbody > tr.full_table")
	rows.Each(func(i int, row *goquery.Selection) {
		name := row.Find("th[data-stat='franch_name']").Text()
		abbr, exists := row.Find("th[data-stat='franch_name'] > a").Attr("href")
		if exists {
			idParts := strings.Split(abbr, "/")
			if len(idParts) > 2 {
				abbr = strings.TrimSuffix(idParts[2], ".html")
				abbr = getCorrectTeamAbbrevation(abbr)
			}

			winlosspct := scraper.GetTDDataStatFloat(row, "win_loss_pct")
			playoffs := scraper.GetTDDataStatInt(row, "years_playoffs")
			divtitles := scraper.GetTDDataStatInt(row, "years_division_champion")
			conftitles := scraper.GetTDDataStatInt(row, "years_conference_champion")
			championships := scraper.GetTDDataStatInt(row, "years_league_champion")

			teams = append(teams, databasestructs.TeamInfo{Name: name, TeamAbbr: abbr, WinLossPct: winlosspct * 100, Playoffs: playoffs, DivisionTitles: divtitles, ConferenceTitles: conftitles, Championships: championships})
		}
	})

	return teams
}

func getCorrectTeamAbbrevation(name string) string {
	switch name {
	case "NOH":
		return "NOP"
	case "CHA":
		return "CHO"
	case "NJN":
		return "BRK"
	}

	return name
}
