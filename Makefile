.PHONY: build test run clean lint

build:
	go build -o bin/terralings cmd/terralings/main.go

test:
	go test -v ./...

lint:
	go vet ./...

clean:
	rm -rf bin/ .terraform/ .cache/
