DATA_GENERATOR=datagenerator


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

dgen:
	@echo "Building the DataGenrator Binary"
	cd ./data_generator && env GOOS=linux CGO_ENABLED=0 go build -a -o ./${DATA_GENERATOR} .

rdgen:
	docker rmi hotelreservation_datagenerator

try: down rdgen dgen up

logs:
	docker logs -f hotelreservation_datagenerator