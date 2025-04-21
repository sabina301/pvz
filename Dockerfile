FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o pvz ./cmd/pvz/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pvz .

EXPOSE 8080
CMD ["./pvz"]