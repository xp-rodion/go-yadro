package domain

import "xkcd/pkg/database"

func ConverterResultToComic(result database.Result) Comic {
	return Comic{
		URL: result.Url,
	}
}
