package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (d *Database) Create() (bool, error) {
	file, err := os.Create(d.Filename)
	if err != nil {
		return false, fmt.Errorf("error creating database file %s: %s", d.Filename, err)
	}
	defer file.Close()
	return true, nil
}

func (d *Database) WriteComic(comic map[string]map[string]interface{}) {
	if len(comic) == 0 {
		return
	}
	fileInfo, _ := os.Stat(d.Filename)

	data := map[string]map[string]interface{}{}
	if fileInfo.Size() != 0 {
		file, _ := ioutil.ReadFile(d.Filename)
		json.Unmarshal(file, &data)
	}
	for k, v := range comic {
		data[k] = v
	}
	file, _ := json.MarshalIndent(data, "", "  ")
	ioutil.WriteFile(d.Filename, file, fileInfo.Mode())
}

func WriteComics(db *Database, comics [][]byte) {
	for _, byteComic := range comics {
		comic, _ := JsonToComic(byteComic)
		comicMap := comic.ComicToMap()
		db.WriteComic(comicMap)
	}
}
