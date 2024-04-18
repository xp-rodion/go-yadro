BINARY = xkcd

all:
	make build

build:

	go build -o $(BINARY) ./cmd/xkcd
clean:
	go clean
	rm -f $(BINARY) database.json cache.json