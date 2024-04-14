package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type Database struct {
	Filename string
}

func (d *Database) Init(filename string) {
	d.Filename = filename
}

func ValidateDatabase(filename string) bool {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return false
	}
	return true
}

// Create создает БД и сразу инциализирует ее ключами (делается для понимания считанных комиксов)
func (d *Database) Create(amountEntry int) (bool, error) {
	file, err := os.Create(d.Filename)
	if err != nil {
		return false, fmt.Errorf("error creating database file %s: %s", d.Filename, err)
	}
	defer file.Close()
	data := make(map[int]map[string]interface{})
	for i := 1; i < amountEntry+1; i++ {
		data[i] = map[string]interface{}{}
	}
	dataJson, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return false, fmt.Errorf("error marshalling data to json: %s", err)
	}
	err = os.WriteFile(d.Filename, dataJson, 0644)
	if err != nil {
		return false, fmt.Errorf("error writing data to file %s: %s", d.Filename, err)
	}
	return true, nil
}

func (d *Database) Get(id int) (Comic, bool) {
	data := make(map[int]map[string]interface{})
	fileInfo, _ := os.Stat(d.Filename)

	if fileInfo.Size() == 0 {
		return Comic{}, false
	}
	file, _ := os.ReadFile(d.Filename)
	err := json.Unmarshal(file, &data)
	if err != nil {
		return Comic{}, false
	}
	content, ok := data[id]
	if !ok {
		return Comic{}, false
	}
	return ComicFromDBEntry(id, content), true
}

func (d *Database) Add(comic Comic) {
	fileInfo, _ := os.Stat(d.Filename)

	data := make(map[int]map[string]interface{})
	if fileInfo.Size() != 0 {
		file, _ := os.ReadFile(d.Filename)
		json.Unmarshal(file, &data)
	}
	id, aboutComic := comic.ToDBEntry()
	data[id] = aboutComic
	file, _ := json.MarshalIndent(data, "", "\t")
	os.WriteFile(d.Filename, file, fileInfo.Mode())
}

func (d *Database) EmptyEntries() []int {
	fileInfo, _ := os.Stat(d.Filename)

	entries := make([]int, 0)
	data := make(map[int]map[string]interface{})
	if fileInfo.Size() == 0 {
		return entries
	}
	file, err := os.ReadFile(d.Filename)
	if err != nil {
		return entries
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return entries
	}

	for key, value := range data {
		if len(value) == 0 {
			entries = append(entries, key)
		}
	}

	return entries
}
