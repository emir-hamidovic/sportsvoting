package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Player struct {
	Name     string
	ID       string
	Points   float64
	Rebounds float64
	Assists  float64
	Steals   float64
	Blocks   float64
}

func parseAllPlayers() ([]Player, error) {
	alphabet := 'a'
	players := []Player{}
	for {
		url := fmt.Sprintf("https://www.basketball-reference.com/players/%s", string(alphabet))
		req, err := setupRequest(url)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %v", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		// Find all rows in the table
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
					// Add player to list of players
					players = append(players, Player{Name: name, ID: id})
				}
			}
		})

		if alphabet += 1; alphabet == '{' {
			break
		}

		// Wait for 4 seconds before sending the next request to try avoiding the rate limit of 20req/min
		time.Sleep(4 * time.Second)
	}

	return players, nil
}

func setupRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set a User-Agent header to impersonate a browser user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	// Set a random Accept-Language header for each request
	langs := []string{"en-US", "en-GB", "fr-FR", "de-DE", "es-ES", "pt-PT", "it-IT", "ja-JP", "ko-KR", "zh-CN"}
	lang := langs[rand.Intn(len(langs))]
	req.Header.Set("Accept-Language", lang)
	return req, nil
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

func parsePlayersIteratively(players []Player) ([]Player, error) {
	fmt.Println("ID\t\t\tName")
	for _, player := range players {
		fmt.Printf("%s\t\t%s\n", player.ID, player.Name)

		url := fmt.Sprintf("https://www.basketball-reference.com/players/%c/%s.html", player.ID[0], player.ID)
		req, err := setupRequest(url)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %v", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		firstRow := doc.Find("#per_game").First().Find("tbody tr").Last()
		currentSeason := getCurrentSeasonFull()

		// for current season, find th[data-stat=season'] and then find the 'a' element's text and try to match the season
		season := strings.TrimSpace(firstRow.Find("th[data-stat='season'] a").Text())
		if season == currentSeason {
			// Extract the relevant stats and convert them to floats
			player.Points, _ = strconv.ParseFloat(strings.TrimSpace(firstRow.Find("td[data-stat='pts_per_g']").Text()), 64)
			player.Rebounds, _ = strconv.ParseFloat(strings.TrimSpace(firstRow.Find("td[data-stat='trb_per_g']").Text()), 64)
			player.Assists, _ = strconv.ParseFloat(strings.TrimSpace(firstRow.Find("td[data-stat='ast_per_g']").Text()), 64)
			player.Blocks, _ = strconv.ParseFloat(strings.TrimSpace(firstRow.Find("td[data-stat='blk_per_g']").Text()), 64)
			player.Steals, _ = strconv.ParseFloat(strings.TrimSpace(firstRow.Find("td[data-stat='stl_per_g']").Text()), 64)
			fmt.Printf("%s stats in the current season: %.1f points, %.1f rebounds, %.1f assists, %.1f blocks, %.1f steals per game\n", player.Name, player.Points, player.Rebounds, player.Assists, player.Blocks, player.Steals)
		}

		// Wait for 4 seconds before sending the next request to try avoiding the rate limit of 20req/min
		time.Sleep(4 * time.Second)
	}

	return players, nil
}

func main() {
	//players, err := parseAllPlayers()
	players, err := parseAllPlayers()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// players, err = parsePlayersIteratively(players)
	_, err = parsePlayersIteratively(players)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
