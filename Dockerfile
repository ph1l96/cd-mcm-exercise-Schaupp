# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./cmd/api

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /api-server .

EXPOSE 8080

ENTRYPOINT ["./api-server"]
