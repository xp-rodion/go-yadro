package search

import (
	"fmt"
	"xkcd/pkg/database"
	"xkcd/pkg/words"
)

func DBRelevantComics(db database.Database, proposal string, count int) {
	comics := db.Entries()
	source := words.Stemmer(proposal)
	for _, comic := range database.DBRelevantComics(source, comics, count) {
		fmt.Printf("Weight: %d; Url: %s\n", comic.Weight, comic.Url)
	}
}

func IndexRelevantComics(db, index database.Database, proposal string, count int) {
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
