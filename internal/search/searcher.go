package search

import (
	"fmt"
	"xkcd/pkg/database"
	"xkcd/pkg/words"
)

func PrintDBRelevantComics(relevantComics database.Results) {
	for _, comic := range relevantComics {
		fmt.Printf("Weight: %d; Url: %s\n", comic.Weight, comic.Url)
	}
}

func DBRelevantComics(db database.JSON, proposal string) database.Results {
	comics := db.Entries()
	source := words.Stemmer(proposal)
	return database.DBRelevantComics(source, comics)
}

func IndexRelevantComics(db, index database.JSON, proposal string, count int) {
	source := words.Stemmer(proposal)
	comics := index.IndexEntries()
	results := database.IndexRelevantComics(source, comics, count)
	ids := make([]int, len(results))
	for _, result := range results {
		ids = append(ids, result.Id)
	}
	for _, comic := range db.Gets(ids) {
		fmt.Printf("Weight: %d; Url: %s\n", getWeight(results, comic.Id), comic.Url)
	}
}

func getWeight(results database.IndexResults, id int) int {
	for _, result := range results {
		if result.Id == id {
			return result.Weight
		}
	}
	return 0
}
