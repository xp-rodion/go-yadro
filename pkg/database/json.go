package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type JSON struct {
	Filename string
}

func (d *JSON) Init(filename string) {
	d.Filename = filename
}

func ValidateDatabase(filename string) bool {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return false
	}
	return true
}

// Create создает БД и сразу инциализирует ее ключами (делается для понимания считанных комиксов)
func (d *JSON) Create(amountEntry int) (bool, error) {
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

func (d *JSON) CreateIndex() (bool, error) {
	file, err := os.Create(d.Filename)
	if err != nil {
		return false, fmt.Errorf("error creating database file %s: %s", d.Filename, err)
	}
	defer file.Close()
	return true, nil
}

func (d *JSON) Get(id int) (Comic, bool) {
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

func (d *JSON) Gets(ids []int) []Comic {
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

func (d *JSON) Add(comic Comic) {
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

func (d *JSON) Adds(comics []Comic) {
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

func (d *JSON) AddsInIndex(comics []Comic) {
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
		keywords := strings.Split(comic.Keywords, ",")
		for _, word := range keywords {
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

func (d *JSON) IndexEntries() map[string][]int {
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

func (d *JSON) Entries() [] Comic {
	fileInfo, _ := os.Stat(d.Filename)

	comics := make([] Comic, 0)
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

func (d *JSON) MapEntries() map[int]interface{} {
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

func (d *JSON) EmptyEntries() []int {
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
