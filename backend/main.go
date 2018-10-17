package main

import (
	"bufio"
	"fmt"
	//"github.com/rapito/go-spotify/spotify"
	//"github.com/zmb3/spotify"
	"github.com/gin-gonic/gin"
	"os"
	//"path"
)

func main() {
	r := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/results", handleGetResults)

	//p := ginprometheus.NewPrometheus("gin")
	//p.Use(r)

	r.Run(":9202") // TODO - get port from config file
}

func handleGetResults(c *gin.Context) {

}

/*func main() {
	playlists := getPlaylistsFromGoogleScrape("site%3Aspotify.com+inurl%3Aplaylist+voxtrot+nujabes")
	for _, item := range playlists {
		fmt.Println(item.ResultDesc)
	}
	//refreshProxyList()

}*/

func getPlaylistsFromGoogleScrape(url string) []GoogleResult {
	res, _ := googleScrape(url, "uk", "en") // TODO: Make proxies configurable

	fmt.Println(res)
	//playlists := make([]string, len(res))

	//for i, item := range res {
	//playlists[i] = path.Base(item.ResultURL)[0:21]
	//}

	return res
}

func getPlaylistsFromGoogleApi() {
	// TODO
}

func getPlaylistsFromFile(url string) []string {
	lines, _ := readLines(url)

	return lines
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
