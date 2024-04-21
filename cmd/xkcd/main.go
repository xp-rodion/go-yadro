package main

import (
	"fmt"
	"xkcd/internal/parse"
	"xkcd/internal/search"
	"xkcd/internal/utils"
)

func main() {
	configPath, proposal, indexSearch := utils.ParseCLIFlags()
	cnf := utils.InitializeConfig(configPath)
	client := utils.InitializeClient(cnf.Url, cnf.CacheFile, 6)
	db := utils.InitializeDB(cnf.Database, client.ComicsCount)
	index := utils.InitializeIndex(cnf.IndexFile)
	parse.ParallelParseComics(client, db, index, cnf.Goroutines)
	fmt.Println("\nРезультат поиска")
	if indexSearch {
		search.IndexRelevantComics(db, index, proposal, 10)
	} else {
		search.DBRelevantComics(db, proposal, 10)
	}
}
