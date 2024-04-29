package parse

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
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

func ParseWorker(ctx context.Context, wg *sync.WaitGroup, queue chan<- xkcd.Entry, client xkcd.Client, entries []int) {
	defer wg.Done()
	for _, idx := range entries {
		select {
		case <-ctx.Done():
			return
		default:
			entry, ok := client.Get(idx)
			if !ok {
				continue
			}
			queue <- entry
		}
	}
}

func ParallelParseComics(client xkcd.Client, db, index database.Database, amountGoroutines int) {
	fmt.Println("Начало парсинга!")
	var wg sync.WaitGroup
	entries := db.EmptyEntries()
	amountEntries := len(entries)
	fmt.Printf("%d комиксов будут спаршены\n", amountEntries)
	queue := make(chan xkcd.Entry, amountEntries)
	goroutineEntries := amountEntries / amountGoroutines
	comics := make([]database.Comic, amountEntries)
	remainder := amountEntries - (goroutineEntries * amountGoroutines) // подсчет остатка
	wg.Add(amountGoroutines)
	// amountEntries - кол-во всех записей, goroutineEntries - кол-во записей на 1 горутину, amountGoroutines - кол-во горутин
	notifyChan := make(chan bool, 1)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	for i := 0; i < amountGoroutines; i++ {
		start := i*goroutineEntries + 1
		end := start + goroutineEntries - 1
		if i == amountGoroutines-1 {
			end += remainder
		}
		go ParseWorker(ctx, &wg, queue, client, entries[start-1:end])
	}

	go func() {
		wg.Wait()
		close(queue)
		notifyChan <- true
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Прерываю программу...")
			stop()
			break Loop
		case <-notifyChan:
			fmt.Println("Комиксы преобразованы, идет запись в бд, прерывание невозможно!")
			break Loop
		}
	}

	fmt.Println("Загрузка в бд...")
	for entry := range queue {
		comic := ConverterEntryToComic(entry)
		comics = append(comics, comic)
	}
	index.AddsInIndex(comics)
	db.Adds(comics)
	fmt.Println("Конец парсинга!")
}
