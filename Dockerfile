# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /todoapp ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /todoapp /app/
COPY --from=builder /app/docs ./docs/
COPY .env .

RUN apk add --no-cache ca-certificates tzdata

EXPOSE 8080
CMD ["/app/todoapp"]