package xkcd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Cache struct {
	filename string
	Count    int       `json:"count"`
	Date     time.Time `json:"date"`
}

func (c *Cache) Validate() bool {
	if _, err := os.Stat(c.filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Cache) Read() bool {
	if !c.Validate() {
		return false
	}
	fileInfo, _ := os.Stat(c.filename)
	if fileInfo.Size() == 0 {
		return false
	}
	file, _ := os.ReadFile(c.filename)
	err := json.Unmarshal(file, c)
	if err != nil {
		return false
	}
	return true
}

func (c *Cache) Write() bool {
	file, err := os.OpenFile(c.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		return false
	}
	defer file.Close()

	data, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Error marshaling cache data:", err)
		return false
	}

	if _, err = file.Write(data); err != nil {
		fmt.Println("Error writing cache data to file:", err)
		return false
	}

	return true
}

func (c *Cache) Update(count int) {
	c.Count = count
	c.Date = time.Now()
	c.Write()
}
