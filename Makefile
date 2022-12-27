.PHONY: true

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-docker:
	docker-compose build todo-app

run-docker:
	docker-compose up todo-app

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up
