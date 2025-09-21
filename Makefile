.PHONY: docker-run

docker-run:
	docker-compose up --build

protoc:
	protoc --go_out=./backend/user-service/internal --go-grpc_out=./backend/user-service/internal  proto/*.proto && \
	protoc --go_out=./backend/auth-service/internal --go-grpc_out=./backend/auth-service/internal  proto/*.proto