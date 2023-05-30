package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetTDDataStatFloat(row *goquery.Selection, id string) float64 {
	res, _ := strconv.ParseFloat(strings.TrimSpace(row.Find(fmt.Sprintf("td[data-stat='%s']", id)).Text()), 64)
	return res
}

func GetTDDataStatInt(row *goquery.Selection, id string) int64 {
	res, _ := strconv.ParseInt(strings.TrimSpace(row.Find(fmt.Sprintf("td[data-stat='%s']", id)).Text()), 10, 64)
	return res
}

func GetTDDataStatString(row *goquery.Selection, id string) string {
	return row.Find(fmt.Sprintf("td[data-stat='%s']", id)).Text()
}

func GetTeamLogo(doc *goquery.Document) string {
	logo, exists := doc.Find("img.teamlogo").Attr("src")
	if exists {
		return logo
	}

	return ""
}
