build:
	go build -o bin/api

run: build
	./bin/api

up:
	docker-compose up -d

down:
	docker-compose down

test:
	go test -v ./...