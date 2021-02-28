build:
	./build.sh

up:
	docker-compose up -d

migrate:
	go run main.go migration up

serve:build
	go run main.go serve-rest

seed:
	node ./migration-script/brand.js