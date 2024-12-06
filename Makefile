.PHONY: build run image start stop

build:
	go build -o bin/gateway ./cmd/main.go

run:
	go run ./cmd/main.go

image:
	docker build -t intro-gateway .

start:
	docker run --name intro-gateway -d --env-file .env intro-gateway

stop:
	docker stop intro-gateway
	docker rm intro-gateway
