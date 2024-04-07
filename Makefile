BINARY = xkcd

all:
	make build

build:
	go build -o $(BINARY) cmd/xkcd/main.go cmd/xkcd/utils.go

clean:
	go clean
	rm -f $(BINARY)