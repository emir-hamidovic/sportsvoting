package teams

import (
	"fmt"
	"log"
	"scraper/database"
	"scraper/parser"
	"scraper/parser/players"
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

func ParseTeams(db database.Database) ([]players.Player, error) {
	allTeams, err := ParseAll()
	if err != nil {
		return nil, err
	}

	var players []players.Player
	var logo string
	for _, team := range allTeams {
		players, logo, err = GetTeamRosterAndLogo(team.TeamAbbr, players)
		if err != nil {
			return nil, err
		}

		team.Logo = logo
		fmt.Println(team)
		_, err = db.InsertTeam(team.TeamAbbr, team.Name, team.Logo, team.WinLossPct, team.Playoffs, team.DivisionTitles, team.ConferenceTitles, team.Championships)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(4 * time.Second)
	}

	fmt.Println("Teams added to database.")
	return players, nil
}

func GetTeamRosterAndLogo(teamAbbr string, players []players.Player) ([]players.Player, string, error) {
	url := fmt.Sprintf("https://www.basketball-reference.com/teams/%s/2023.html", teamAbbr)
	res, err := parser.SendRequest(url)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", err
	}

	players = append(players, findRosterInfo(doc, teamAbbr)...)
	return players, getTeamLogo(doc), nil
}

func getTeamLogo(doc *goquery.Document) string {
	logo, exists := doc.Find("img.teamlogo").Attr("src")
	if exists {
		return logo
	}

	return ""
}

func findRosterInfo(doc *goquery.Document, team string) []players.Player {
	rows := doc.Find("table#roster > tbody > tr")
	var playersList []players.Player
	rows.Each(func(i int, row *goquery.Selection) {
		name := row.Find("td[data-stat='player']").Text()
		fmt.Println(name)
		id, exists := row.Find("td[data-stat='player'] > a").Attr("href")
		if exists {
			idParts := strings.Split(id, "/")
			if len(idParts) > 3 {
				id = strings.TrimSuffix(idParts[3], ".html")
			}

			fmt.Printf("%s: %s\n", id, name)

			college := row.Find("td[data-stat='college']").Last().Text()
			height := row.Find("td[data-stat='height']").Text()
			weight := row.Find("td[data-stat='weight']").Text()

			playersList = append(playersList, players.Player{Name: name, ID: id, College: college, Height: height, Weight: weight, TeamAbbr: team})
		}
	})

	return playersList
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

	return findBasicTeamInfo(doc, teams), nil
}

func findBasicTeamInfo(doc *goquery.Document, teams []Teams) []Teams {
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

			winlosspct, _ := strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='win_loss_pct']").Text()), 64)
			playoffs, _ := strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='years_playoffs']").Text()), 10, 64)
			divtitles, _ := strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='years_division_champion']").Text()), 10, 64)
			conftitles, _ := strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='years_conference_champion']").Text()), 10, 64)
			championships, _ := strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='years_league_champion']").Text()), 10, 64)

			teams = append(teams, Teams{Name: name, TeamAbbr: abbr, WinLossPct: winlosspct * 100, Playoffs: playoffs, DivisionTitles: divtitles, ConferenceTitles: conftitles, Championships: championships})
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
