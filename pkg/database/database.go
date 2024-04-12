package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

func (d *Database) Create() (bool, error) {
	file, err := os.Create(d.Filename)
	if err != nil {
		return false, fmt.Errorf("error creating database file %s: %s", d.Filename, err)
	}
	defer file.Close()
	return true, nil
}

func (d *Database) Get(id int) (Comic, bool) {
	Id := strconv.Itoa(id)
	data := make(map[string]map[string]interface{})
	fileInfo, _ := os.Stat(d.Filename)
	if fileInfo.Size() == 0 {
		return Comic{}, false
	}
	file, _ := os.ReadFile(d.Filename)
	err := json.Unmarshal(file, &data)
	if err != nil {
		return Comic{}, false
	}
	content, ok := data[Id]
	if !ok {
		return Comic{}, false
	}
	return ComicFromDBEntry(Id, content), true
}

func (d *Database) Add(comic Comic) {
	fileInfo, _ := os.Stat(d.Filename)

	data := make(map[string]map[string]interface{})
	if fileInfo.Size() != 0 {
		file, _ := os.ReadFile(d.Filename)
		json.Unmarshal(file, &data)
	}
	id, aboutComic := comic.ToDBEntry()
	data[id] = aboutComic
	file, _ := json.MarshalIndent(data, "", "\t")
	os.WriteFile(d.Filename, file, fileInfo.Mode())
}
