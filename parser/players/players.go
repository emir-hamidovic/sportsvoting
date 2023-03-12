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

	players, err := parseIteratively(allPlayers)
	if err != nil {
		return err
	}

	for _, player := range players {
		_, err = db.InsertPlayer(player.ID, player.Name, player.College, player.TeamAbbr, player.Height, player.Weight, player.Age, player.Games, player.Minutes, player.Points, player.Rebounds, player.Assists, player.Steals, player.Blocks, player.Turnovers, player.FGPercentage*100, player.FTPercentage*100, player.ThreeFGPercentage*100)
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

func parseIteratively(players []Player) ([]Player, error) {
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

		currentSeasonRow := doc.Find("#per_game").First().Find("tbody tr.full_table").Last()
		var currentTeam string
		if currentSeasonRow.Find("td[data-stat='team_id']").Text() == "TOT" {
			currentTeam = doc.Find("#per_game").First().Find("tbody tr").Last().Find("td[data-stat='team_id']").Text()
		}

		players[i] = fillPlayerInfoForSeason(currentSeasonRow, player, currentTeam)
		fmt.Printf("%s stats in the current season: %s ID, %s college, %s teamabbr, %s height, %s weight, %d age, %d games, %.1f minutes, %.1f points, %.1f rebounds, %.1f assists, %.1f blocks, %.1f steals %.1f turnovers, %.1f fg, %.1f ft, %.1f 3pt\n", player.Name, player.ID, players[i].College, players[i].TeamAbbr, players[i].Height, players[i].Weight, players[i].Age, players[i].Games, players[i].Minutes, players[i].Points, players[i].Rebounds, players[i].Assists, players[i].Blocks, players[i].Steals, players[i].Turnovers, players[i].FGPercentage*100, players[i].FTPercentage*100, players[i].ThreeFGPercentage*100)

		// Wait for 5 seconds before sending the next request to try avoiding the rate limit of 20req/min
		time.Sleep(5 * time.Second)
	}

	return players, nil
}

func fillPlayerInfoForSeason(row *goquery.Selection, player Player, currentTeam string) Player {
	if currentTeam == "" {
		player.TeamAbbr = getOldTeamAbbreviation(strings.TrimSpace(row.Find("td[data-stat='team_id']").Text()))
	} else {
		player.TeamAbbr = getOldTeamAbbreviation(currentTeam)
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

func getOldTeamAbbreviation(abbr string) string {
	switch abbr {
	case "NOP":
		return "NOH"
	case "CHO":
		return "CHA"
	case "BRK":
		return "NJN"
	}

	return abbr
}

func ParseBoxScores(db database.Database) error {
	url := "https://www.basketball-reference.com/boxscores"
	res, err := parser.SendRequest(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	divs := doc.Find("div.game_summaries > div.game_summary")
	divs.Each(func(i int, row *goquery.Selection) {
		var players []Player
		winningTeam := row.Find("table.teams > tbody > tr.winner")
		players, err := appendPlayersFromTodaysPlayedGames(players, winningTeam, db)
		if err != nil {
			return
		}

		losingTeam := row.Find("table.teams > tbody > tr.loser")
		players, err = appendPlayersFromTodaysPlayedGames(players, losingTeam, db)
		if err != nil {
			return
		}

		players, err = parseIteratively(players)
		if err != nil {
			return
		}

		for _, val := range players {
			_, err := db.UpdatePlayers(val.Games, val.Minutes, val.Points, val.Rebounds, val.Assists, val.Blocks, val.Steals, val.Turnovers, val.FGPercentage, val.FTPercentage, val.ThreeFGPercentage, val.ID)
			if err != nil {
				return
			}
		}
	})

	return nil
}

func appendPlayersFromTodaysPlayedGames(players []Player, selection *goquery.Selection, db database.Database) ([]Player, error) {
	abbr, exists := selection.Find("td > a").First().Attr("href")
	if exists {
		idParts := strings.Split(abbr, "/")
		if len(idParts) > 3 {
			abbr = getOldTeamAbbreviation(idParts[2])
			rows, err := db.SelectPlayerID(abbr)
			if err != nil {
				return nil, err
			}

			for rows.Next() {
				var player Player
				err := rows.Scan(&player.ID)
				if err != nil {
					return nil, err
				}

				players = append(players, player)
			}
		}
	}

	return players, nil
}
