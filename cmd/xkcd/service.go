package main

import (
	"fmt"
	"sync"
	"xkcd/pkg/database"
	"xkcd/pkg/xkcd"
)

// ParseComics парс одного комикса и его запись в БД (для промежуточной записи)
func ParseComics(client xkcd.Client, db database.Database) {
	entries := db.EmptyEntries()
	for _, idx := range entries {
		entry, ok := client.Get(idx)
		if !ok {
			continue
		}
		comic := ConverterEntryToComic(entry)
		db.Add(comic)
	}
	fmt.Print("Parse ended!")
}

func ParseWorker(wg *sync.WaitGroup, queue chan<- xkcd.Entry, client xkcd.Client, entries []int) {
	defer wg.Done()
	for _, idx := range entries {
		entry, ok := client.Get(idx)
		if !ok {
			continue
		}
		queue <- entry
	}
}

func ParallelParseComics(client xkcd.Client, db database.Database, amountGoroutines int) {
	fmt.Println("Начало парсинга!")
	var mtx sync.Mutex
	var wg sync.WaitGroup
	entries := db.EmptyEntries()
	amountEntries := len(entries)
	fmt.Printf("%d комиксов будут спаршены\n", amountEntries)
	queue := make(chan xkcd.Entry, amountEntries)
	// amountEntries - кол-во всех записей, goroutineEntries - кол-во записей на 1 горутину, amountGoroutines - кол-во горутин
	goroutineEntries := amountEntries / amountGoroutines
	for i := 0; i < amountGoroutines; i++ {
		wg.Add(1)
		fmt.Printf("%d/%d горутин учавствует в парсинге\n", i+1, amountGoroutines)
		start := i*goroutineEntries + 1
		end := start + goroutineEntries - 1
		go ParseWorker(&wg, queue, client, entries[start-1:end])
	}

	go func() {
		wg.Wait()
		close(queue)
		fmt.Println("Загрузка в бд...")
	}()

	for entry := range queue {
		comic := ConverterEntryToComic(entry)
		mtx.Lock()
		db.Add(comic)
		mtx.Unlock()
	}
	fmt.Println("Конец парсинга!")
}
