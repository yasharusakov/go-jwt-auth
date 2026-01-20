FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG APP_PORT

RUN go build -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app .

EXPOSE ${APP_PORT}

CMD ["./app"]