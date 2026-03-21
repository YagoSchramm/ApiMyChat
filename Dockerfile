### Multi-stage build for the Go API
# Stage 1: builder (match go.mod >=1.24)
FROM golang:1.25-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux

# Download deps first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN go build -o api ./cmd/api

# Stage 2: minimal runtime image
FROM gcr.io/distroless/static-debian12 AS runtime
WORKDIR /app

COPY --from=builder /app/api .

# The binary reads API_PORT (defaults to 8000 in code)
ENV API_PORT=8000

EXPOSE 8000
CMD ["./api"]
