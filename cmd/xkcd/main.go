package main

import (
	"xkcd/pkg/config"
	"xkcd/pkg/database"
)

func main() {
	configPath, amountComics, loggingToConsole := parseCLIFlags()
	cnf := new(config.Config)
	cnf.Init(configPath)
	client := initializeClient(cnf.Url, cnf.ClientLogFile)
	db := initializeDB(cnf.Database)
	byteComics := client.GetFixedComics(amountComics, loggingToConsole)
	database.WriteComics(db, byteComics)
}
