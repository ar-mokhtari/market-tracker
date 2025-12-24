# Change from 1.21 to 1.23
FROM golang:1.23-alpine AS builder

WORKDIR /app

# The rest of your build steps...
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# Copy config if needed
# COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./main"]
