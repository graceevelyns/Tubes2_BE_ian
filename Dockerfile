FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
RUN go mod verify
COPY src/ ./src/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /app/alchemy-backend ./src/cmd/main.go

# runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/alchemy-backend /app/alchemy-backend

COPY src/cmd/docs ./docs

# port server
EXPOSE 8080

# run binary backend
CMD ["/app/alchemy-backend"]