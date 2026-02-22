# Gunakan base image golang untuk build
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go.mod dan go.sum dulu untuk caching dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Stage kedua: image ringan untuk menjalankan binary
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates untuk koneksi TLS
RUN apk add --no-cache ca-certificates

# Copy binary dari stage builder
COPY --from=builder /app/main .

# Copy folder web, uploads, storage, dll
COPY ./web ./web
COPY ./uploads ./uploads
COPY ./storage ./storage
COPY .env .env

# Expose port
EXPOSE 3000

# Jalankan aplikasi
CMD ["./main"]
