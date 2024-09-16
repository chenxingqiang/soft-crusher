# Build stage
FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o soft-crusher cmd/soft-crusher/main.go

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/soft-crusher .
COPY --from=builder /app/config/config.yaml .
EXPOSE 8080
CMD ["./soft-crusher"]