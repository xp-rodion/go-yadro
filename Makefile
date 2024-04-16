BINARY = xkcd

all:
	make build

build:
	go build -o $(BINARY) cmd/xkcd/main.go cmd/xkcd/utils.go cmd/xkcd/converter.go cmd/xkcd/service.go

clean:
	go clean
	rm -f $(BINARY) database.json cache.json