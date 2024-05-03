package utils

import (
	"flag"
	"fmt"
	"xkcd/pkg/config"
	"xkcd/pkg/database"
	"xkcd/pkg/words"
	"xkcd/pkg/xkcd"
)

func InitializeDB(filename string, amountEntry int) database.Database {
	db := database.Database{}
	db.Init(filename)
	status := database.ValidateDatabase(db.Filename)
	if status {
		db.Create(amountEntry)
	}
	return db
}

func InitializeIndex(filename string) database.Database {
	db := database.Database{}
	db.Init(filename)
	status := database.ValidateDatabase(db.Filename)
	if status {
		db.CreateIndex()
	}
	return db
}

func InitializeConfig(filename string) config.Config {
	cnf := config.Config{}
	cnf.Init(filename)
	return cnf
}

func InitializeClient(url, logFile string, timeout int) xkcd.Client {
	client := xkcd.Client{}
	client.Init(url, "info.0.json", logFile, timeout)
	return client
}

func ParseCLIFlags() (configFile, proposal string, indexSearch bool, port string) {
	flag.StringVar(&configFile, "c", "configs/config.yaml", "Configuration file")
	flag.BoolVar(&indexSearch, "i", false, "Index Search")
	flag.StringVar(&port, "p", "", "Port")
	proposal = words.ReadCLIArgs()

	flag.Parse()

	if len(configFile) == 0 {
		fmt.Println("Don't see Configuration file, use default Configuration file")
	}

	if indexSearch {
		fmt.Println("Using Index Search")
	} else {
		fmt.Println("Using Default Search")
	}

	return
}

func CheckNewComics(db database.Database, client xkcd.Client) []int {
	entries := db.MapEntries()
	newComics := make([]int, 0)
	for i := 1; i < client.ComicsCount+1; i++ {
		if i == 404 {
			continue
		}
		if _, suc := entries[i]; suc != true {
			newComics = append(newComics, i)
		}
	}
	return newComics
}
