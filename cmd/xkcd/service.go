package main

import (
	"fmt"
	"xkcd/pkg/database"
	"xkcd/pkg/xkcd"
)

// ParseComics парс одного комикса и его запись в БД (для промежуточной записи)
func ParseComics(client xkcd.Client, db database.Database, amount int) {
	for idx := 1; idx < client.ComicsCount+1; idx++ {
		entry, ok := client.Get(idx)
		if !ok {
			continue
		}
		comic := ConverterEntryToComic(entry)
		if amount > 0 {
			fmt.Println(comic)
			amount--
		}
		db.Add(comic)
	}
	fmt.Print("Parse ended!")
}
