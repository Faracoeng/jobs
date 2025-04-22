# Makefile

test:
	go test ./... -v

build-api: test
	docker build -f build/Dockerfile.api -t solucao-api .

build-etl: test
	docker build -f build/Dockerfile.etl -t solucao-etl .

build-generator:
	docker build -f build/Dockerfile.generator -t csv-generator .

build-all: test build-api build-etl build-generator

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart: down up

test-cover:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out

run: build-all up
