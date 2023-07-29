build:
	@go build -o bin/app

run:
	@./bin/app server --port=3001 --datastore

test:
	go test -v ./... -count=1

vendor:
	go mod vendor

lint:
	golangci-lint run -v ./...
