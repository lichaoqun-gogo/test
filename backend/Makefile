.PHONY: run run-dev run-prod build clean install-deps

run:
	go run cmd/server/main.go

run-dev:
	APP_ENV=dev go run cmd/server/main.go

run-prod:
	APP_ENV=prod go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

clean:
	rm -rf bin/

install-deps:
	go mod tidy
