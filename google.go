package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// Represets a single search result
type GoogleResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

// In case we want to use different domains...
var googleDomains = map[string]string{
	"com": "https://www.google.com/search?q=",
	"uk":  "https://www.google.co.uk/search?q=",
	"ru":  "https://www.google.ru/search?q=",
	"fr":  "https://www.google.fr/search?q=",
}

// In case we want to rotate user agents...
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

// Helper function to get a random user agent
func randUserAgent() string {
	rand.Seed(time.Now().Unix())
	return userAgents[rand.Intn(len(userAgents))]
}

// Helper function to build a search query
func buildGoogleUrl(searchTerm string, countryCode string, languageCode string) string {
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found {
		return fmt.Sprintf("%s%s&num=100&hl=%s", googleBase, searchTerm, languageCode)
	} else {
		return fmt.Sprintf("%s%s&num=100&hl=%s", googleDomains["com"], searchTerm, languageCode)
	}
}

// Helper function to attempt to make a HTTP request
// TODO: Detect if requests are being throttled and request new IP address
func googleRequest(searchURL string) (*http.Response, error) {

	baseClient := &http.Client{}

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randUserAgent())

	res, err := baseClient.Do(req)

	log.Info("Created request with status: " + res.Status)

	if err != nil {
		log.Error("Error creating request: " + err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

// Helper functin to parse a response and return an array of GoogleResults
func googleResultParser(response *http.Response) ([]GoogleResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)

	file, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	str, _ := doc.Html()
	fmt.Fprintf(file, str)

	if err != nil {
		return nil, err
	}
	results := []GoogleResult{}
	sel := doc.Find("div.g")
	rank := 1
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")
		if link != "" && link != "#" {
			result := GoogleResult{
				rank,
				link,
				title,
				desc,
			}
			results = append(results, result)
			rank += 1
		}
	}
	return results, err
}

// Given search parameters, make a request and return an array of GoogleResults
func GoogleScrape(searchTerm string, countryCode string, languageCode string) ([]GoogleResult, error) {
	googleUrl := buildGoogleUrl(searchTerm, countryCode, languageCode)
	res, err := googleRequest(googleUrl)
	if err != nil {
		return nil, err
	}

	scrapes, err := googleResultParser(res)
	if err != nil {
		return nil, err
	} else {
		return scrapes, nil
	}
}
