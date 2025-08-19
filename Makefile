.PHONY: server client docker-run

server:
	cd server && go run cmd/app/main.go

client:
	cd client && npm run dev

docker-run:
	docker-compose up --build