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

try: down rdgen dgen up

logs:
	docker logs -f hotelreservation_datagenerator

.PHONY: proto

proto:
	$(PROTOC) \
		--go_out=$(GO_OUT_DIR) \
		--go-grpc_out=$(GO_OUT_DIR) \
		--proto_path=$(PROTO_DIR) \
		$(PROTO_DIR)/*.proto