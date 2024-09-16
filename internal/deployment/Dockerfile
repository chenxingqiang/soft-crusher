# File: Dockerfile

# Build the Go binary
FROM golang:1.17 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o soft-crusher ./cmd/soft-crusher

# Build the React frontend
FROM node:14 as frontend-builder
WORKDIR /app
COPY frontend .
RUN npm install && npm run build

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/soft-crusher .
COPY --from=frontend-builder /app/build ./frontend/build
COPY config.yaml .
EXPOSE 8080
CMD ["./soft-crusher"]