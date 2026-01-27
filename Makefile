.PHONY: docker-run swagger swagger-auth

docker-run:
	docker-compose up --build

swagger: swagger-auth
	@echo "Swagger docs generated for all services"

swagger-auth:
	@echo "Generating Swagger for auth-service..."
	cd ./backend/auth-service && swag init -g cmd/main.go -o docs

#swagger-user:
#	@echo "Generating Swagger for user-service..."
#	cd ./backend/user-service && swag init -g cmd/main.go -o docs

#swagger-gateway:
#	@echo "Generating Swagger for api-gateway..."
#	cd ./backend/api-gateway && swag init -g cmd/main.go -o docs