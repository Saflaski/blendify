FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal


RUN go build -o app ./cmd/api

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .
COPY .env .env

EXPOSE 3000
CMD ["./app", "-logtostderr=true"]
