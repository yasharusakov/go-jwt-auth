.PHONY: server client

server:
	cd server && go run cmd/app/main.go

client:
	cd client && npm run dev