package main

import (
	"xkcd/pkg/database"
	"xkcd/pkg/xkcd"
)

// ParseComics парс одного комикса и его запись в БД (для промежуточной записи)
func ParseComics(client xkcd.Client, db database.Database) {
	for idx := 1; idx < client.ComicsCount+1; idx++ {
		entry, ok := client.Get(idx)
		if !ok {
			continue
		}
		comic := ConverterEntryToComic(entry)
		db.Add(comic)
	}
}
