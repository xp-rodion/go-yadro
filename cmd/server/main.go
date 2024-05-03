package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xkcd/internal/core/services"
	"xkcd/internal/handlers"
	"xkcd/internal/parse"
	"xkcd/internal/repositories"
	"xkcd/internal/server"
	"xkcd/internal/utils"
)

func main() {
	configPath, port, _, _ := utils.ParseCLIFlags()
	cnf := utils.InitializeConfig(configPath)
	client := utils.InitializeClient(cnf.Url, cnf.CacheFile, 6)
	db := utils.InitializeDB(cnf.Database, client.ComicsCount)
	entries := db.EmptyEntries()
	comics := parse.ParallelParseComics(client, entries, cnf.Goroutines)
	db.Adds(comics)

	repo := repositories.NewJsonRepository(db, client)
	service := services.NewService(repo)
	handler := handlers.NewHTTPHandler(service)
	if port == "" {
		port = cnf.Port
	}

	go func() {
		duration := time.Hour * 24
		for {
			fmt.Println("update run every 24 hours")
			_, err := repo.Update()
			if err != nil {
				log.Fatal("error occurred while running update func", err)
			}
			time.Sleep(duration)
		}
	}()

	srv := new(server.Server)
	go func() {
		if err := srv.Run(port, *handler); err != nil {
			log.Fatal("error occurred while starting server", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("error occured on server shutting down", err.Error())
	}

}
