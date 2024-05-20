package parse

import (
	"strings"
	"xkcd/pkg/database"
	"xkcd/pkg/words"
	"xkcd/pkg/xkcd"
)

// ConverterEntryToComic Конвертер из Client-представления (при парсе) в Database-представление (для записи в БД)
func ConverterEntryToComic(entry xkcd.Entry) database.Comic {
	keywords := strings.Join(append(words.Stemmer(entry.Alt), words.Stemmer(entry.Title)...), ",")
	return database.Comic{
		Id:       entry.Id,
		Url:      entry.Url,
		Keywords: keywords,
	}
}
