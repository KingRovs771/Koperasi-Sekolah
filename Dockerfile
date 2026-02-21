# Tahap 1: Build Aplikasi
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy go.mod dan go.sum lalu download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o kopsis-app main.go

# Tahap 2: Runner (Image super ringan)
FROM alpine:latest
WORKDIR /app

# Install zona waktu dan client PostgreSQL (untuk fitur Backup & Restore Anda)
RUN apk --no-cache add ca-certificates tzdata postgresql-client
ENV TZ=Asia/Jakarta

# Copy file hasil build dari tahap 1
COPY --from=builder /app/kopsis-app .

# Copy folder statis dan views (sesuai struktur Anda: folder 'web')
COPY --from=builder /app/web ./web

# Buat folder penyimpanan agar tidak error saat upload/backup
RUN mkdir -p /app/storage/backups /app/uploads

EXPOSE 8080

CMD ["./kopsis-app"]