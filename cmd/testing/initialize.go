package main

import (
	"xkcd/internal/parse"
	"xkcd/internal/utils"
)

func main() {
	configPath := "configs/benchmark_config.yaml"
	cnf := utils.InitializeConfig(configPath)
	client := utils.InitializeClient(cnf.Url, cnf.CacheFile, 6)
	db := utils.InitializeDB(cnf.Database, client.ComicsCount)
	index := utils.InitializeIndex(cnf.IndexFile)
	entries := db.EmptyEntries()
	comics := parse.ParallelParseComics(client, entries, cnf.Goroutines)
	db.Adds(comics)
	index.AddsInIndex(comics)
}
