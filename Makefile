build:
	@go build -o bin/app

run:
	@./bin/app server

test:
	go test -v ./... -count=1