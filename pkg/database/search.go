package database

import (
	"sort"
)

func DBRelevantComics(source []string, entries []Comic, count int) Results {
	sourceMap := make(map[string]bool, len(source))
	for _, entry := range source {
		sourceMap[entry] = true
	}
	results := Search(sourceMap, entries)
	sort.Sort(results)
	return results[:count]
}

func IndexRelevantComics(source []string, entries map[string][]int, count int) IndexResults {
	results := make(map[int]int)
	for _, word := range source {
		ids := entries[word]
		for _, id := range ids {
			results[id] = results[id] + 1
		}
	}
	idxResult := MapToIndexResult(results)
	sort.Sort(idxResult)
	if count > idxResult.Len() {
		count = idxResult.Len()
	}
	return idxResult[:count]
}

func Search(source map[string]bool, entries []Comic) Results {
	results := make(Results, 0, len(entries))
	for _, entry := range entries {
		results = append(results, Result{Weight: Weight(source, entry), Url: entry.Url, Keywords: entry.Keywords})
	}
	return results
}

func Weight(source map[string]bool, entry Comic) int {
	keywords := entry.Keywords
	weight := 0
	for _, word := range keywords {
		_, ok := source[word]
		if ok {
			weight++
		}
	}
	return weight
}
