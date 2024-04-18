package main

import (
	"flag"
	"fmt"
	"xkcd/pkg/config"
	"xkcd/pkg/database"
	"xkcd/pkg/xkcd"
)

func initializeDB(filename string, amountEntry int) database.Database {
	db := database.Database{}
	db.Init(filename)
	status := database.ValidateDatabase(db.Filename)
	if status {
		db.Create(amountEntry)
	}
	return db
}

func initializeConfig(filename string) config.Config {
	cnf := config.Config{}
	cnf.Init(filename)
	return cnf
}

func initializeClient(url, logFile, cacheFile string, amountGoroutines, timeout int) xkcd.Client {
	client := xkcd.Client{}
	client.Init(url, "info.0.json", logFile, cacheFile, amountGoroutines, timeout)
	return client
}

func parseCLIFlags() (configFile string) {
	flag.StringVar(&configFile, "c", "configs/config.yaml", "Configuration file")
	flag.Parse()
	if len(configFile) == 0 {
		fmt.Println("Don't see Configuration file, use default Configuration file")
	}

	return
}
