package main

import (
	"fmt"
	"xkcd/internal/parse"
	"xkcd/internal/search"
	"xkcd/internal/utils"
)

func main() {
	configPath, proposal, indexSearch, _ := utils.ParseCLIFlags()
	cnf := utils.InitializeConfig(configPath)
	client := utils.InitializeClient(cnf.Url, cnf.CacheFile, 6)
	db := utils.InitializeJSON(cnf.Database, client.ComicsCount)
	index := utils.InitializeIndex(cnf.IndexFile)
	entries := db.EmptyEntries()
	comics := parse.ParallelParseComics(client, entries, cnf.Goroutines)
	db.Adds(comics)
	index.AddsInIndex(comics)
	fmt.Println("\nРезультат поиска")
	if indexSearch {
		search.IndexRelevantComics(db, index, proposal, 10)
	} else {
		comics := search.DBRelevantComics(db, proposal)
		search.PrintDBRelevantComics(comics[:10])
	}
}
