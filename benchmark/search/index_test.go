package search

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"xkcd/internal/parse"
	"xkcd/internal/search"
	"xkcd/internal/utils"
	"xkcd/pkg/database"
)

func initializeIdx() (db, index database.Database) {
	configPath := "../configs/config.yaml"
	cnf := utils.InitializeConfig(configPath)
	client := utils.InitializeClient(cnf.Url, cnf.CacheFile, 6)
	db = utils.InitializeDB(cnf.Database, client.ComicsCount)
	index = utils.InitializeIndex(cnf.IndexFile)
	parse.ParallelParseComics(client, db, index, cnf.Goroutines)
	return
}

func BenchmarkIndexSearch(b *testing.B) {
	proposal := "I'm following on your questions"

	fmt.Println("INDEX SEARCH")

	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	db, index := initializeIdx()
	search.IndexRelevantComics(db, index, proposal, 10)

	w.Close()
	os.Stdout = stdout
	io.Copy(ioutil.Discard, r)
}
