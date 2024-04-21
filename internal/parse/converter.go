package parse

import (
	"xkcd/pkg/database"
	"xkcd/pkg/words"
	"xkcd/pkg/xkcd"
)

// ConverterEntryToComic Конвертер из Client-представления (при парсе) в Database-представление (для записи в БД)
func ConverterEntryToComic(entry xkcd.Entry) database.Comic {
	return database.Comic{
		Id:       entry.Id,
		Url:      entry.Url,
		Keywords: words.Stemmer(entry.Alt),
	}
}
