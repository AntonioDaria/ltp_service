DIRS=$(shell find packages/ -mindepth 1 -type d | awk '{gsub("//","/"); printf "%s,", $$0}' | sed 's/,$$//')

swag:
	@swag i -d api/,$(DIRS) -g server.go -o docs/static/api

start:
	@docker-compose up -d --build

stop:
	@docker-compose down