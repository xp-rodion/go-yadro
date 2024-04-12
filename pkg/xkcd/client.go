package xkcd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Client struct {
	BaseUrl     string
	Client      *http.Client
	ComicsCount int
	UrlFormat   string
	logFile     string
}

func (c *Client) Init(url string, format string, logFile string, timeout int) {
	c.BaseUrl = url
	c.Client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
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
	err = json.NewDecoder(resp.Body).Decode(&data)
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

func (c *Client) Get(id int) (Entry, bool) {
	url := c.reverse(id)
	client := c.Client
	entry := Entry{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return entry, false
	}
	resp, err := client.Do(req)
	if err != nil && c.validateStatusCode(resp.StatusCode) {
		return entry, false
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&entry)
	if err != nil {
		return Entry{}, false
	}
	return entry, true
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

func (c *Client) validateStatusCode(statusCode int) bool {
	return statusCode >= 500 && statusCode <= 505
}
