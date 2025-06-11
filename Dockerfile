FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/url-shortener-api

FROM alpine:latest AS app
WORKDIR /app
EXPOSE 8080
CMD ["./server"]
COPY --from=builder /app/server .
