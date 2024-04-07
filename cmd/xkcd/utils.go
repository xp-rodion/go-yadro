package main

import (
	"flag"
	"fmt"
	"xkcd/pkg/database"
	"xkcd/pkg/xkcd"
)

func initializeDB(filename string) *database.Database {
	db := new(database.Database)
	db.Init(filename)
	status := database.ValidateDatabase(db.Filename)
	if status {
		db.Create()
	}
	return db
}

func initializeClient(url string, logFile string) *xkcd.Client {
	client := new(xkcd.Client)
	client.Init(url, "info.0.json", logFile)
	return client
}

func parseCLIFlags() (configFile string, amountComics int, loggingToConsole bool) {

	flag.StringVar(&configFile, "c", "configs/config.yaml", "Configuration file")
	flag.IntVar(&amountComics, "n", 0, "Number of comics to download")
	flag.BoolVar(&loggingToConsole, "o", false, "Enable console logging")

	flag.Parse()

	if len(configFile) == 0 {
		fmt.Println("Don't see Configuration file, use default Configuration file\n")
	}

	if amountComics == 0 {
		fmt.Println("Parse all comics!\n")
	}

	if loggingToConsole {
		fmt.Println("Logging enabled\n")
	}

	return
}
