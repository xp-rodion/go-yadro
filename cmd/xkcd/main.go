package main

import (
	"xkcd/pkg/config"
)

func main() {
	configPath, amount, logging := parseCLIFlags()
	cnf := new(config.Config)
	cnf.Init(configPath)
	client := initializeClient(cnf.Url, cnf.ClientLogFile, 5)
	db := initializeDB(cnf.Database)
	if !logging {
		amount = 0
	}
	ParseComics(client, db, amount)
}
