package main

import (
	"fmt"
	"github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	//"github.com/corpix/uarand"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GoogleResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

var googleDomains = map[string]string{
	"com": "https://www.google.com/search?q=",
	"uk":  "https://www.google.co.uk/search?q=",
	"ru":  "https://www.google.ru/search?q=",
	"fr":  "https://www.google.fr/search?q=",
}

var proxyList []string

func refreshProxyList() {
	baseClient := &http.Client{}

	req, _ := http.NewRequest("GET", "https://www.sslproxies.org/", nil)
	req.Header.Set("User-Agent", browser.Random())

	res, _ := baseClient.Do(req)
	doc, _ := goquery.NewDocumentFromResponse(res)

	//rows := 0
	// Better way to parse this html?
	doc.Find("#proxylisttable").Find("tr").Each(func(_ int, tr *goquery.Selection) {

		proxy := "http://" + tr.Find("td").Eq(0).Text() + ":" + tr.Find("td").Eq(1).Text()
		if len(proxy) > 9 { // Check if IP is valid (TODO: find better way)
			proxyList = append(proxyList, proxy)
		}
	})
}

func buildGoogleUrl(searchTerm string, countryCode string, languageCode string) string {
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found {
		return fmt.Sprintf("%s%s&num=100&hl=%s", googleBase, searchTerm, languageCode)
	} else {
		return fmt.Sprintf("%s%s&num=100&hl=%s", googleDomains["com"], searchTerm, languageCode)
	}
}

func googleRequest(searchURL string) (*http.Response, error) {

	if len(proxyList) == 0 {
		refreshProxyList()
	}

	//proxyString := randomString(proxyList)
	//fmt.Println(proxyString)
	proxyUrl, err := url.Parse("//69.75.136.106:42045")

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	req, err := http.NewRequest("GET", searchURL, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return res, nil
	}
}

func googleResultParser(response *http.Response) ([]GoogleResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
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
		titleTag := item.Find("h3.r")
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

func googleScrape(searchTerm string, countryCode string, languageCode string) ([]GoogleResult, error) {
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

func randomString(options []string) string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(options)
	return options[randNum]
}
