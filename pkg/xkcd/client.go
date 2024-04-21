package xkcd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type Client struct {
	BaseUrl     string
	Client      *http.Client
	ComicsCount int
	UrlFormat   string
	logFile     string
	Cache       Cache
}

func (c *Client) Init(url, format, logFile string, timeout int) {
	c.BaseUrl = url
	c.Client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	c.ComicsCount = OldGetComicsCount(url, format)
	c.UrlFormat = format
	c.logFile = logFile
}

func (c *Client) CacheFileName() {
	fmt.Println("Cache:", c.Cache.filename)
}

func OldGetComicsCount(url string, format string) int {
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

func (c *Client) Search(amountGoroutines int) int {
	fmt.Println("Ищу последний комикс...")
	wg := sync.WaitGroup{}
	wg.Add(amountGoroutines)
	queue := make(chan int, amountGoroutines)
	step := 10
	k := step * amountGoroutines
	for i := 0; i < amountGoroutines; i++ {
		go c.SearchWorker(&wg, queue, i, step, k)
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	ids := make([]int, 0)
	for id := range queue {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	fmt.Println("Последний комикс найден!")
	return ids[0] - 1
}

func (c *Client) getComicsCount(amountGoroutines int) int {
	cache := c.Cache
	if !cache.Validate() {
		comicsCount := c.Search(amountGoroutines)
		cache.Update(comicsCount)
		return comicsCount
	}
	cache.Read()
	if cache.Date.Add(time.Minute * 10).After(time.Now()) {
		return cache.Count
	}
	comicsCount := c.Search(amountGoroutines)
	cache.Update(comicsCount)
	return comicsCount
}

func (c *Client) SearchWorker(wg *sync.WaitGroup, queue chan<- int, nWorker, step, k int) {
	defer wg.Done()
	iter := 0
	client := c.Client
	for {
		start := (k * iter) + (nWorker * step) + 1
		end := start + step
		for idx := start; idx < end; idx++ {
			url := c.reverse(idx)
			resp, err := client.Get(url)
			if err != nil {
				continue
			}
			status := resp.StatusCode
			if status == 404 && idx != 404 {
				queue <- idx
				return
			}
		}
		iter++
	}
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
		return entry, false
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
