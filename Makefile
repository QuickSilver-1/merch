.PHONY: up down

up:
	docker build . --file Dockerfile -t app:latest
	docker-compose up

down:
	docker-compose down