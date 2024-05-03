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
		if i == 404 {
			continue
		}
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

func (d *Database) CreateIndex() (bool, error) {
	file, err := os.Create(d.Filename)
	if err != nil {
		return false, fmt.Errorf("error creating database file %s: %s", d.Filename, err)
	}
	defer file.Close()
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

func (d *Database) Gets(ids []int) []Comic {
	comics := make([]Comic, 0)
	data := make(map[int]map[string]interface{})
	fileInfo, _ := os.Stat(d.Filename)

	if fileInfo.Size() == 0 {
		return comics
	}
	file, _ := os.ReadFile(d.Filename)
	err := json.Unmarshal(file, &data)
	if err != nil {
		return comics
	}
	for _, id := range ids {
		content, ok := data[id]
		if !ok {
			continue
		}
		comics = append(comics, ComicFromDBEntry(id, content))
	}
	return comics
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

func (d *Database) Adds(comics []Comic) {
	fileInfo, _ := os.Stat(d.Filename)

	data := make(map[int]map[string]interface{})
	if fileInfo.Size() != 0 {
		file, _ := os.ReadFile(d.Filename)
		json.Unmarshal(file, &data)
	}
	for _, comic := range comics {
		if comic.Id == 0 {
			continue
		}
		id, aboutComic := comic.ToDBEntry()
		data[id] = aboutComic
	}
	file, _ := json.MarshalIndent(data, "", "\t")
	os.WriteFile(d.Filename, file, fileInfo.Mode())
}

func (d *Database) AddsInIndex(comics []Comic) {
	fileInfo, _ := os.Stat(d.Filename)
	data := make(map[string][]int)
	if fileInfo.Size() != 0 {
		file, _ := os.ReadFile(d.Filename)
		json.Unmarshal(file, &data)
	}
	for _, comic := range comics {
		if comic.Id == 0 {
			continue
		}
		for _, word := range comic.Keywords {
			ids, ok := data[word]
			if ok {
				data[word] = append(ids, comic.Id)
			} else {
				data[word] = []int{comic.Id}
			}
		}
	}
	file, _ := json.MarshalIndent(data, "", "\t")
	os.WriteFile(d.Filename, file, fileInfo.Mode())
}

func (d *Database) IndexEntries() map[string][]int {
	fileInfo, _ := os.Stat(d.Filename)
	data := make(map[string][]int)
	if fileInfo.Size() == 0 {
		return data
	}
	file, err := os.ReadFile(d.Filename)
	if err != nil {
		return data
	}
	json.Unmarshal(file, &data)
	return data
}

func (d *Database) Entries() []Comic {
	fileInfo, _ := os.Stat(d.Filename)

	comics := make([]Comic, 0)
	data := make(map[int]map[string]interface{})
	if fileInfo.Size() == 0 {
		return comics
	}
	file, err := os.ReadFile(d.Filename)
	if err != nil {
		return comics
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return comics
	}
	for key, value := range data {
		if value["keywords"] == nil {
			continue
		}
		comics = append(comics, ComicFromDBEntry(key, value))
	}
	return comics
}

func (d *Database) MapEntries() map[int]interface{} {
	fileInfo, _ := os.Stat(d.Filename)
	entries := make(map[int]interface{})
	if fileInfo.Size() == 0 {
		return entries
	}
	file, err := os.ReadFile(d.Filename)
	if err != nil {
		return entries
	}
	err = json.Unmarshal(file, &entries)
	if err != nil {
		return entries
	}
	return entries
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
