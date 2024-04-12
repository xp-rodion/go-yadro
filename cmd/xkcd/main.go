package main

import (
	"xkcd/pkg/config"
)

func main() {
	configPath, _, _ := parseCLIFlags()
	cnf := new(config.Config)
	cnf.Init(configPath)
	client := initializeClient(cnf.Url, cnf.ClientLogFile, 1)
	db := initializeDB(cnf.Database)
	ParseComics(client, db)
}
