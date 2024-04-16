package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	configPath := parseCLIFlags()
	cnf := initializeConfig(configPath)
	client := initializeClient(cnf.Url, cnf.ClientLogFile, cnf.CacheFile, cnf.Goroutines, 6)
	db := initializeDB(cnf.Database, client.ComicsCount)
	ParallelParseComics(client, db, cnf.Goroutines)
	end := time.Now()
	fmt.Println("Время выполнения:", end.Sub(start))
}
