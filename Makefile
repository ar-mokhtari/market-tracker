include .env

start: db-up run

db-up:
	docker compose up -d

run:
	go run main.go

fetch:
	curl -X POST http://localhost:8080/api/v1/prices/fetch

prices:
	curl http://localhost:8080/api/v1/prices

gold:
	curl -X GET http://localhost:8080/api/v1/prices?type=gold

currency:
	curl -X GET http://localhost:8080/api/v1/prices?type=currency
