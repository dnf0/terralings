.PHONY: build test run clean lint

build:
	go build -o bin/terralings ./cmd/terralings

run:
	go run ./cmd/terralings

test:
	go test -v -race ./...

lint:
	go vet ./...

clean:
	rm -rf bin/ .terraform/ .cache/
