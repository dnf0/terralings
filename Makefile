.PHONY: all build build-all run test test-race lint fmt-check verify check clean

all: check build build-all test-race

build:
	go build -o bin/terralings ./cmd/terralings

build-all:
	go build -v ./...

run:
	go run ./cmd/terralings

test:
	go test -v ./...

test-race:
	go test -v -race ./...

lint:
	go vet ./...

fmt-check:
	@unformatted=$$(gofmt -l .) || exit 1; \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted files:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

verify:
	go mod verify

check: verify fmt-check lint

clean:
	rm -rf bin/ terralings .terraform/ .cache/
