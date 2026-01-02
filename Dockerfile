# Build Stage
FROM golang:1.23-alpine AS builder

# Gerekli sistem paketlerini yükle (gcc, musl-dev cgo için gerekebilir)
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Dependency'leri önbellekle
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodu kopyala
COPY . .

# Uygulamayı derle (Static Binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o pulseguard cmd/pulseguard/main.go

# Final Stage (Sadece binary'i al, kaynak kodu sil)
FROM alpine:latest

WORKDIR /root/

# SSL sertifikaları ve Timezone verisi lazım
RUN apk --no-cache add ca-certificates tzdata

# Builder'dan binary'i kopyala
COPY --from=builder /app/pulseguard .
# Frontend dosyalarını kopyala (HTML/JS/CSS)
COPY --from=builder /app/web ./web
# Config dosyasını kopyala (fallback için)
COPY --from=builder /app/config.yaml ./config.yaml

# Portu aç
EXPOSE 8080

# Çalıştır
CMD ["./pulseguard"]
