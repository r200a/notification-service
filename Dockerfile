FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o notification-server ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/notification-server .
COPY .env .
CMD ["./notification-server"]