package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func main() {
	r := gin.Default()

	var port string
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.GET("/results", handleGetResults)
	r.GET("/health", handleGetHealth)
	if os.Getenv("PORT") == "" {
		port = "9207"
	} else {
		port = os.Getenv("PORT")
	}
	r.Run(":" + port) // TODO - get port from config file
}

type Result struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func handleGetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, "I'm alive")
}

func handleGetResults(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

	q := c.Request.URL.Query()
	list := q["list"][0]
	terms := strings.Split(list, ",")
	queryString := "site%3Aspotify.com+inurl%3Aplaylist+"

	for i, term := range terms {
		fmt.Println(term)
		queryString += "\"" + strings.Replace(term, " ", "+", -1)  + "\""
		if i != len(terms)-1{
			queryString += "+"
		}
	}
	
	results := getPlaylistsFromGoogleScrape(queryString)
	parsedResults := make([]Result, len(results))

	for i, result := range results {
		parsedResults[i].Title = result.ResultTitle
		fmt.Println(result.ResultTitle)

		// It seems that based on the browser, the url string may contain q=
		urlThing := strings.Split(result.ResultURL, "q=")
		if len(urlThing) > 1 {
			parsedResults[i].Url = urlThing[1]
		} else {
			parsedResults[i].Url = urlThing[0]
		}
	}

	c.JSON(http.StatusOK, parsedResults)
}

func getPlaylistsFromGoogleScrape(url string) []GoogleResult {
	res, _ := GoogleScrape(url, "com", "en")
	return res
}
