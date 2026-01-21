FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG APP_PORT

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add curl

COPY --from=builder /app/app .

EXPOSE ${APP_PORT}

CMD ["./app"]