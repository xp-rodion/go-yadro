package xkcd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Client struct {
	BaseUrl     string
	Client      *http.Client
	ComicsCount int
	UrlFormat   string
	logFile     string
}

func (c *Client) Init(url string, format string, logFile string) {
	c.BaseUrl = url
	c.Client = &http.Client{}
	c.ComicsCount = getComicsCount(url, format)
	c.UrlFormat = format
	c.logFile = logFile
}

func getComicsCount(url string, format string) int {
	url = fmt.Sprintf("%s/%s", url, format)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	count := data["num"].(float64)
	return int(count)
}

func (c *Client) reverse(id int) string {
	url := fmt.Sprintf("%s/%d/%s", c.BaseUrl, id, c.UrlFormat)
	return url
}

func (c *Client) getComic(id int) ([]byte, error) {
	url := c.reverse(id)
	client := c.Client
	body := make([]byte, 0)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return body, err
	}
	resp, err := client.Do(req)
	if err != nil && validateStatusCode(resp.StatusCode) {
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func (c *Client) GetFixedComics(amount int, loggingToConsole bool) (comics [][]byte) {

	badRequests := make([]int, 0)

	if amount > c.ComicsCount || amount < 1 {
		amount = c.ComicsCount
	}
	step := float64(amount) / 100.0

	for i := 1; i < amount+1; i++ {
		comic, err := c.getComic(i)
		if err != nil {
			badRequests = append(badRequests, i)
			continue
		}
		if loggingToConsole {
			fmt.Println(string(comic), "\n\n")
		} else {
			fmt.Println(int(float64(i)/step), "%")
		}

		comics = append(comics, comic)
	}
	fmt.Println("Writed!")
	c.LoggingBadRequest(badRequests)
	return comics
}

func (c *Client) GetComicsByIDs(ids []int) (comics [][]byte, err error) {
	for _, id := range ids {
		comic, err := c.getComic(id)
		if err != nil {
			return nil, err
		}
		comics = append(comics, comic)
	}
	return comics, nil
}

//LoggingBadRequest логгирование неудачных запросов (будет сохраняться в файле), их можно будет скормить GetComicsByIDs
func (c *Client) LoggingBadRequest(comics []int) {

	if len(comics) == 0 {
		return
	}

	logFile, _ := os.OpenFile(c.logFile, os.O_CREATE, 0644)
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	for _, comic := range comics {
		logger.Printf("%d,", comic)
	}
}

func validateStatusCode(statusCode int) bool {
	return statusCode >= 500 && statusCode <= 505
}
