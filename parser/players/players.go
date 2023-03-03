package players

import (
	"fmt"
	"scraper/database"
	"scraper/parser"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Player struct {
	Name              string
	ID                string
	College           string
	TeamAbbr          string
	Height            string
	Weight            string
	Age               int64
	Games             int64
	Minutes           float64
	Points            float64
	Rebounds          float64
	Assists           float64
	Steals            float64
	Blocks            float64
	Turnovers         float64
	FGPercentage      float64
	ThreeFGPercentage float64
	FTPercentage      float64
}

func ParsePlayers(db database.Database) error {
	allPlayers, err := ParseAll()
	if err != nil {
		return err
	}
	fmt.Println("Players parsed")

	players, err := ParseIteratively(allPlayers)
	if err != nil {
		return err
	}

	stmt, err := db.PrepareStatementForPlayerInsert()
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, player := range players {
		_, err = stmt.Exec(player.ID, player.Name, player.College, player.TeamAbbr, player.Height, player.Weight, player.Age, player.Games, player.Minutes, player.Points, player.Rebounds, player.Assists, player.Steals, player.Blocks, player.Turnovers, player.FGPercentage*100, player.FTPercentage*100, player.ThreeFGPercentage*100)
		if err != nil {
			return err
		}
	}

	fmt.Println("Players added to database.")
	return nil
}

func ParseAll() ([]Player, error) {
	alphabet := 'a'
	players := []Player{}
	for {
		url := fmt.Sprintf("https://www.basketball-reference.com/players/%s", string(alphabet))
		res, err := parser.SendRequest(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		players = findBasicPlayerInfo(doc, players)

		if alphabet += 1; alphabet == '{' {
			break
		}

		// Wait for 5 seconds before sending the next request to try avoiding the rate limit of 20req/min
		time.Sleep(5 * time.Second)
	}

	return players, nil
}

func findBasicPlayerInfo(doc *goquery.Document, players []Player) []Player {
	rows := doc.Find("table#players > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		toYearStr := row.Find("td[data-stat='year_max']").Text()
		if toYearStr == getEndYearOfTheSeason() {
			name := row.Find("th[data-stat='player']").Text()
			id, exists := row.Find("th[data-stat='player'] > strong > a").Attr("href")
			if exists {
				// Extract player ID from URL
				idParts := strings.Split(id, "/")
				if len(idParts) > 3 {
					id = strings.TrimSuffix(idParts[3], ".html")
				}

				fmt.Printf("%s: %s\n", id, name)

				college := row.Find("td[data-stat='colleges']").Last().Text()
				height := row.Find("td[data-stat='height']").Text()
				weight := row.Find("td[data-stat='weight']").Text()

				players = append(players, Player{Name: name, ID: id, College: college, Height: height, Weight: weight})
			}
		}
	})

	return players
}

func getEndYearOfTheSeason() string {
	today := time.Now()
	year := today.Year()
	month := int(today.Month())
	var currentSeason string
	if month < 10 {
		currentSeason = fmt.Sprint(year)
	} else {
		currentSeason = fmt.Sprint(year + 1)
	}

	return currentSeason
}

func ParseIteratively(players []Player) ([]Player, error) {
	for i, player := range players {
		url := fmt.Sprintf("https://www.basketball-reference.com/players/%c/%s.html", player.ID[0], player.ID)
		res, err := parser.SendRequest(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		firstRow := doc.Find("#per_game").First().Find("tbody tr").Last()
		players[i] = fillPlayerInfoForSeason(firstRow, player)
		fmt.Printf("%s stats in the current season: %s ID, %s college, %s teamabbr %s height %s weight %d age %d games %.1f minutes %.1f points, %.1f rebounds, %.1f assists, %.1f blocks, %.1f steals %.1f turnovers %.1f fg %.1f ft %.1f 3pt per game\n", player.Name, player.ID, players[i].College, players[i].TeamAbbr, players[i].Height, players[i].Weight, players[i].Age, players[i].Games, players[i].Minutes, players[i].Points, players[i].Rebounds, players[i].Assists, players[i].Blocks, players[i].Steals, players[i].Turnovers, players[i].FGPercentage*100, players[i].FTPercentage*100, players[i].ThreeFGPercentage*100)

		// Wait for 5 seconds before sending the next request to try avoiding the rate limit of 20req/min
		time.Sleep(5 * time.Second)
	}

	return players, nil
}

/*
func isPlayerIsPlayingInTheCurrentSeason(row *goquery.Selection) bool {
	currentSeason := getCurrentSeasonFull()
	season := strings.TrimSpace(row.Find("th[data-stat='season'] a").Text())
	return season == currentSeason
}

func getCurrentSeasonFull() string {
	today := time.Now()
	year := today.Year()
	month := int(today.Month())
	var currentSeason string
	if month < 10 {
		currentSeason = fmt.Sprintf("%d-%d", year-1, year%100)
	} else {
		currentSeason = fmt.Sprintf("%d-%d", year, (year+1)%100)
	}

	return currentSeason
}*/

func fillPlayerInfoForSeason(row *goquery.Selection, player Player) Player {
	player.TeamAbbr = strings.TrimSpace(row.Find("td[data-stat='team_id']").Text())
	switch player.TeamAbbr {
	case "NOP":
		player.TeamAbbr = "NOH"
	case "CHO":
		player.TeamAbbr = "CHA"
	case "BRK":
		player.TeamAbbr = "NJN"
	}

	player.Games, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='g']").Text()), 10, 64)
	player.Age, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='age']").Text()), 10, 64)
	player.Minutes, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='mp_per_g']").Text()), 64)
	player.Points, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='pts_per_g']").Text()), 64)
	player.Rebounds, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='trb_per_g']").Text()), 64)
	player.Assists, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ast_per_g']").Text()), 64)
	player.Blocks, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='blk_per_g']").Text()), 64)
	player.Steals, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='stl_per_g']").Text()), 64)
	player.Turnovers, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='tov_per_g']").Text()), 64)
	player.FGPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='fg_pct']").Text()), 64)
	player.ThreeFGPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='fg3_pct']").Text()), 64)
	player.FTPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ft_pct']").Text()), 64)

	return player
}
