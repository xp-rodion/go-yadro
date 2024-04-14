package main

func main() {
	configPath := parseCLIFlags()
	cnf := initializeConfig(configPath)
	client := initializeClient(cnf.Url, cnf.ClientLogFile, 6)
	db := initializeDB(cnf.Database, client.ComicsCount)
	ParallelParseComics(client, db, cnf.Goroutines)
}
