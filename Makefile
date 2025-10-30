.PHONY: docker-run protoc sql-up-auth sql-down-auth sql-up-user sql-down-user

docker-run:
	docker-compose up --build

# migrate create -ext sql -dir ./backend/auth-service/migrations -seq create_refresh_tokens
# migrate create -ext sql -dir ./backend/user-service/migrations -seq create_users

sql-up-auth:
	migrate -path ./backend/auth-service/migrations -database "postgres://postgres:admin@localhost:5488/auth_db?sslmode=disable" up

# connect to the auth_db
# docker-compose exec postgres-auth psql -U postgres -d auth_db

sql-down-auth:
	migrate -path ./backend/auth-service/migrations -database "postgres://postgres:admin@localhost:5488/auth_db?sslmode=disable" down

sql-up-user:
	migrate -path ./backend/user-service/migrations -database "postgres://postgres:admin@localhost:5499/user_db?sslmode=disable" up

# connect to the user_db
# docker-compose exec postgres-user psql -U postgres -d user_db

sql-down-user:
	migrate -path ./backend/user-service/migrations -database "postgres://postgres:admin@localhost:5499/user_db?sslmode=disable" down

protoc:
	protoc --go_out=./backend/user-service/internal --go-grpc_out=./backend/user-service/internal  proto/*.proto && \
	protoc --go_out=./backend/auth-service/internal --go-grpc_out=./backend/auth-service/internal  proto/*.proto