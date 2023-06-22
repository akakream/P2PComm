build:
	@go build -o bin/app

run:
	@./bin/app server --port=3000 --servertype=libp2p --datastore

test:
	go test -v ./... -count=1

vendor:
	@go mod vendor
