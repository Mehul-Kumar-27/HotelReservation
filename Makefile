DATA_GENERATOR=datagenerator
# Specify the path to the protoc executable
PROTOC := $(shell which protoc)

# Specify the path to the directory containing the proto files
PROTO_DIR := ./proto

# Specify the output directory for the generated Go files
GO_OUT_DIR := ./proto/gen


build:
	go build -o bin/api

run: build
	./bin/api

up:
	docker-compose up -d

down:
	docker-compose down

test:
	go test -v -race -count=1 ./... ;

dgen:
	@echo "Building the DataGenrator Binary"
	cd ./data_generator && env GOOS=linux CGO_ENABLED=0 go build -a -o ./${DATA_GENERATOR} .

rdgen:
	docker rmi hotelreservation_datagenerator

try: down up

logs:
	docker logs -f hotelreservation_datagenerator

.PHONY: proto

proto:
	$(PROTOC) \
		--go_out=$(GO_OUT_DIR) \
		--go-grpc_out=$(GO_OUT_DIR) \
		--proto_path=$(PROTO_DIR) \
		$(PROTO_DIR)/*.proto

# http://localhost:8083/connectors
# {
#     "name": "reservation-connector",  
#   "config": {  
#     "connector.class": "io.debezium.connector.mysql.MySqlConnector",
#     "tasks.max": "1",  
#     "database.hostname": "mysql",  
#     "database.port": "3306",
#     "database.user": "mehul",
#     "database.password": "mehulpassword",
#     "database.server.id": "1",  
#     "database.server.name": "reservationserver1",  
#     "database.include.list": "reservation",  
#     "database.history.kafka.bootstrap.servers": "broker:29092",  
#     "database.history.kafka.topic": "schema-changes.inventory" 
# }