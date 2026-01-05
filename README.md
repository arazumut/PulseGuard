# PulseGuard 🛡️

**PulseGuard**, performans ve ölçeklenebilirliğin ön planda tutulduğu, **Go (Golang)** ile geliştirilmiş, kendi sunucunuzda barındırabileceğiniz (self-hosted) modern bir servis izleme ve anomali tespit sistemidir.


## 🏗 Mimari Yaklaşım (Architecture)

Bu proje, kodun test edilebilirliğini, bakımını ve ölçeklenebilirliğini sağlamak amacıyla **Hexagonal Architecture (Ports and Adapters)** ilkelerine sadık kalınarak tasarlanmıştır.

### Katmanlar

1.  **Core (Domain Layer) `internal/core`**:
    *   Uygulamanın kalbidir. İş kuralları (Business Logic) ve Entity'ler burada bulunur.
    *   *Dış dünyadan (DB, HTTP, Redis) habersizdir.* Framework bağımsızdır.
    *   `ports` paketi, dış dünya ile iletişim kurmak için gerekli `interface` tanımlarını içerir.

2.  **Adapters (Infrastructure Layer) `internal/adapter`**:
    *   Core katmanındaki portları (interface) implemente eder.
    *   **Handler**: HTTP isteklerini karşılar (`Fiber` web framework).
    *   **Storage**: Veritabanı ve Cache işlemlerini yapar (`PostgreSQL`, `Redis`).

3.  **Monitor Engine `internal/monitor`**:
    *   Sistemin "Motor" kısmıdır.
    *   Binlerce servisi aynı anda kontrol etmek için **Worker Pool** pattern kullanır.
    *   Non-blocking G/Ç için Go Concurrency (Goroutines & Channels) yoğun olarak kullanılır.

### 🛠 Teknoloji Yığını

*   **Dil**: Go 1.25+
*   **Web Framework**: Fiber (Hız ve düşük bellek tüketimi için)
*   **Database**: PostgreSQL (Kalıcı veri ve zaman serisi benzeri yapılar için)
*   **Cache/Queue**: Redis (Anlık durum yönetimi ve job kuyruğu için)
*   **Logging**: `slog` (Structured Logging - JSON formatında)
*   **Config**: `viper` (Environment variable yönetimi)

## 🚀 Kurulum (Development)

Sistem **Docker-First** yaklaşımıyla tasarlanmıştır.

```bash
# Bağımlılıkları yükle
go mod tidy

# Projeyi çalıştır (Local)
go run cmd/pulseguard/main.go
```

## 📡 API Kullanımı

Sistem şu an In-Memory çalışmaktadır. Uygulama kapandığında veriler silinir.

### 1. Servis Ekleme

Monitoring başlatmak için servisi kaydedin:

```bash
curl -X POST http://localhost:8080/api/v1/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Google",
    "url": "https://google.com",
    "interval": 10
  }'
```

### 2. Servisleri Listeleme

```bash
curl http://localhost:8080/api/v1/services
```

### 3. Logları İzleme

Uygulama loglarında her 10 saniyede bir `Health Check` çıktısı göreceksiniz:

```json
{"time":"...","level":"INFO","msg":"Health Check","service":"Google","status_code":200,"latency":123456789,"success":true}
```

## 🗺️ Roadmap Durumu

- [x] **Faz 0:** Mimari Kurulum (Hexagonal + Fiber)
- [x] **Faz 1:** Core Monitoring Engine (Scheduler + Pinger)
- [x] **Faz 2:** Akıllı Analiz (Latency trendleri)
- [x] **Faz 3:** PostgreSQL & Redis (Kalıcı veri)
