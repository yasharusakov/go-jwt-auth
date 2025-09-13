.PHONY: api-gateway auth-service user-service client docker-run

api-gateway:
	cd backend/api-gateway && go run cmd/main.go

api-service:
	cd backend/auth-service && go run cmd/main.go

user-service:
	cd backend/user-service && go run cmd/main.go

client:
	cd frontend/client && npm run dev

docker-run:
	docker-compose up --build