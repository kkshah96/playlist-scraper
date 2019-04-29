package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sort"
	"strings"
)

func main() {
	r := gin.Default()

	var port string

  // Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.GET("/results", handleGetResults)
	r.GET("/health", handleGetHealth)

  // Heroku will set a PORT env variable, so we will default to using this if
  // it exists
  if os.Getenv("PORT") == "" {
		port = "9207"
	} else {
		port = os.Getenv("PORT")
	}
	r.Run(":" + port) // TODO - get port from config file
}

// Represents a single Spotify search result
type Result struct {
	Title     string         `json:"title"`
	Url       string         `json:"url"`
	Hits      map[string]int `json:"hits"`
	TotalHits int            `json:"toatlHits"`
}

// Health endpoint for testing purpsoes
func handleGetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, "I'm alive!")
}

// Functionality for making a request to Google and processing the results into
// Result struct format
func handleGetResults(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

	q := c.Request.URL.Query()
	fmt.Println(q["list"][0])

	listArr := q["list"]
	doRank := len(q["rank"]) > 0 && q["rank"][0] == "true"

	var list string

	if len(listArr) > 0 {
		list = listArr[0]
	} else {
		c.JSON(http.StatusBadRequest, "Must provide list of query parameters")
		return
	}

	terms := strings.Split(list, ",")
	queryString := "site%3Aspotify.com+inurl%3Aplaylist+"

	for i, term := range terms {
		fmt.Println(term)
		queryString += "\"" + strings.Replace(term, " ", "+", -1) + "\""
		if i != len(terms)-1 {
			queryString += "+"
		}
	}

	results := getPlaylistsFromGoogleScrape(queryString)
	parsedResults := make([]Result, len(results))

	for i, result := range results {
		elementMap := make(map[string]int)

		parsedResults[i].Title = result.ResultTitle
		fmt.Println(result.ResultTitle)

		// It seems that based on the browser, the url string may contain q=
		urlThing := strings.Split(result.ResultURL, "q=")
		if len(urlThing) > 1 {
			parsedResults[i].Url = urlThing[1]
		} else {
			parsedResults[i].Url = urlThing[0]
		}

		if doRank == true {

			for _, element := range terms {
				fmt.Println(element)
				fmt.Println("-------")
				fmt.Println(result.ResultDesc)

				if strings.Contains(strings.ToUpper(result.ResultDesc), strings.ToUpper(element)) {
					elementMap[element]++
					parsedResults[i].TotalHits++
				}

			}

			parsedResults[i].Hits = elementMap
		}
	}

	if doRank == true {
		// Sort first by number of keywords found, then
		// by number of total hits
		// https://stackoverflow.com/questions/36122668/golang-how-to-sort-struct-with-multiple-sort-parameters
		sort.Slice(parsedResults, func(i, j int) bool {
			if len(parsedResults[i].Hits) > len(parsedResults[j].Hits) {
				return true
			}

			if len(parsedResults[i].Hits) < len(parsedResults[j].Hits) {
				return false
			}

			return parsedResults[i].TotalHits > parsedResults[j].TotalHits
		})
	}

	c.JSON(http.StatusOK, parsedResults)
}

// Helper function to call GoogleScrape to fetch GoogleResults
func getPlaylistsFromGoogleScrape(url string) []GoogleResult {
	res, _ := GoogleScrape(url, "com", "en")
	return res
}
