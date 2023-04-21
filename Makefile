.DEFAULT_GOAL = build-prod

LIB = dist/pinion

go.mod:
	go mod tidy
	go mod verify
.PHONY: go.mod

build:
	go build -o $(LIB) ./
.PHONY: build

build-prod: go.mod
	go clean -cache; go build -o $(LIB) -ldflags="-s -w" ./
.PHONY: build-prod

test:
	go test -cover ./...
.PHONY: test

test-prod: go.mod
	go clean -cache; go test -cover ./...
.PHONY: test
