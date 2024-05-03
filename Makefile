BINARY = xkcd
BINARY_SERVER = xkcd-server

all:
	make build

build:
	go build -o $(BINARY) ./cmd/xkcd

benchmark_without_db:
	 rm -rf ./internal/testing/*.json && go test -bench=. ./benchmark/search/default_test.go && rm -rf ./internal/testing/*.json && go test -bench=. ./benchmark/search/index_test.go && rm -rf ./internal/testing/*.json

benchmark_with_db:
	 rm -rf ./internal/testing/*.json && go run cmd/testing/initialize.go && go test -bench=. ./benchmark/search/default_test.go && go test -bench=. ./benchmark/search/index_test.go

server:
	go build -o $(BINARY_SERVER) ./cmd/server

clean:
	go clean
	rm -f $(BINARY) $(BINARY_SERVER) *.json