# ==========================================
# STAGE 1: Build biner Golang
# ==========================================
FROM golang:1.25-alpine AS builder

# Install git, ca-certificates, and templ CLI (diperlukan untuk kirim email/webhook HTTPS)
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

# Copy file dependency dulu agar memanfaatkan caching Docker
COPY go.mod go.sum ./
RUN go mod download

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy seluruh source code (termasuk folder assets)
COPY . .

# Generate templ components, then compile Go menjadi biner tunggal
RUN templ generate && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

# ==========================================
# STAGE 2: Runtime Environment
# ==========================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy biner hasil compile dari Stage 1
COPY --from=builder /app/main .

# COPY FOLDER STATIS (assets only — templates are compiled into the binary)
COPY --from=builder /app/assets ./assets

# Ekspos port sesuai aplikasi Gin (misal 8080)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]