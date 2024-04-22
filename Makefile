BINARY = xkcd

all:
	make build

build:
	go build -o $(BINARY) ./cmd/xkcd

benchmark_without_db:
	 rm -rf ./internal/testing/*.json && go test -bench=. ./internal/testing/default_test.go && rm -rf ./internal/testing/*.json && go test -bench=. ./internal/testing/index_test.go && rm -rf ./internal/testing/*.json

benchmark_with_db:
	 rm -rf ./internal/testing/*.json && go run cmd/xkcd/initialize.go && go test -bench=. ./internal/testing/default_test.go && go test -bench=. ./internal/testing/index_test.go


clean:
	go clean
	rm -f $(BINARY) *.json