# PULSEGUARD

> Hafif, hızlı, gerçek zamanlı servis sağlık & anomali izleme sistemi

> Hafif, hızlı, gerçek zamanlı servis sağlık & anomali izleme sistemi

---

## 🎯 PROJE HEDEFİ

Küçük–orta ölçekli ekiplerin **ağır monitoring sistemleri kurmadan** servislerini izleyebilmesini sağlamak.

* Down olmadan önce sorun tespiti
* Latency / error artışı erken uyarı
* Self-hosted veya SaaS çalışabilen mimari
* **Tamamen Docker tabanlı kurulum (Backend + PostgreSQL + Redis)**
* **Hazır iconic UI template entegrasyonu**

---

# 🧭 FAZ 0 – PROJE TEMELİ & PLANLAMA (Sprint 0)

### 🎯 Amaç

Sağlam temel, net scope, teknik borçsuz başlangıç

### 🔹 Task List

* [ ] Ürün ismi & branding netleştirme (**PulseGuard**)
* [ ] Problem tanımı & hedef kullanıcı profili yazımı
* [ ] Rakip analizi (UptimeRobot, Datadog, NewRelic)
* [ ] MVP kapsamının netleştirilmesi
* [ ] Teknoloji stack kararı (Go, PostgreSQL, Redis)
* [ ] **Docker-first mimari kararı**
* [ ] Repo oluşturma (monorepo)
* [ ] Coding standartları ve branch stratejisi

### 📦 Çıktılar

* README.md (vizyon + hedef)
* Architecture overview diyagramı
* Docker Compose taslak dosyası

---

# 🚀 FAZ 1 – CORE ENGINE (Sprint 1)

### 🎯 Amaç

Sistemin kalbi: **yüksek performanslı heartbeat motoru**

### 🔹 Task List

* [ ] Go proje yapısının oluşturulması
* [ ] HTTP server (net/http veya Fiber)
* [ ] `/heartbeat` endpoint
* [ ] Request validation
* [ ] In-memory service registry
* [ ] Goroutine bazlı heartbeat worker
* [ ] Timeout detection logic
* [ ] İlk logging altyapısı

### 📦 Çıktılar

* Çalışan heartbeat API
* 1000+ servis simülasyonu testi

---

# 🧠 FAZ 2 – AKILLI SAĞLIK ANALİZİ (Sprint 2)

### 🎯 Amaç

Sadece DOWN değil, **sorun yaklaşırken fark etmek**

### 🔹 Task List

* [ ] Latency trend hesaplama
* [ ] Error rate moving average
* [ ] Threshold bazlı uyarı sistemi
* [ ] Basit anomaly detection (z-score)
* [ ] Service state machine (Healthy / Warning / Critical)
* [ ] State transition kuralları

### 📦 Çıktılar

* Sağlık skorlaması
* Anomali tespit edilen servisler

---

# 🗄️ FAZ 3 – VERİ KATMANI (Sprint 3)

### 🎯 Amaç

Kalıcı, hızlı ve ölçeklenebilir veri altyapısı

### 🔹 Task List

* [ ] PostgreSQL schema tasarımı
* [ ] Service table
* [ ] Metrics table (time-series light)
* [ ] Redis cache entegrasyonu
* [ ] Data retention policy
* [ ] Repository pattern implementasyonu

### 📦 Çıktılar

* Kalıcı metric kayıtları
* Performanslı veri okuma

---

# 📡 FAZ 4 – REAL-TIME DASHBOARD API (Sprint 4)

### 🎯 Amaç

Frontend’e canlı veri sağlayan API katmanı

### 🔹 Task List

* [ ] Service list endpoint
* [ ] Service detail endpoint
* [ ] Health status endpoint
* [ ] WebSocket / SSE altyapısı
* [ ] Real-time push mekanizması
* [ ] Pagination & filtering

### 📦 Çıktılar

* Dashboard için hazır API
* Canlı servis güncellemeleri

---

# 🎨 FAZ 5 – DASHBOARD & UX (Sprint 5)

### 🎯 Amaç

Hazır **iconic UI template** üzerine hızlı ve temiz entegrasyon

### 🔹 Task List

* [ ] Iconic template proje yapısına entegrasyon
* [ ] Genel durum ekranı (overview)
* [ ] Servis sağlık kartları
* [ ] Latency & error grafik ekranları
* [ ] Warning / Critical renk sistemi
* [ ] Real-time veri binding (WebSocket / SSE)
* [ ] Responsive düzenlemeler

### 📦 Çıktılar

* Profesyonel dashboard
* Canlı güncellenen UI

---

# 🔔 FAZ 6 – ALERT & NOTIFICATION (Sprint 6)

### 🎯 Amaç

Sorun olduğunda **anında haber vermek**

### 🔹 Task List

* [ ] Alert rule engine
* [ ] Email notification
* [ ] Slack webhook entegrasyonu
* [ ] Alert throttling
* [ ] Alert history kaydı

### 📦 Çıktılar

* Çalışan uyarı sistemi
* Alert geçmişi

---

# 🐳 FAZ 7 – DEPLOY & SELF-HOSTED (Sprint 7)

### 🎯 Amaç

Tek komutla ayağa kalkabilen **tam Docker stack**

### 🔹 Task List

* [ ] Go backend için multi-stage Dockerfile
* [ ] PostgreSQL Docker container
* [ ] Redis Docker container
* [ ] Docker Compose (prod & local)
* [ ] Env config yapısı (.env)
* [ ] Healthcheck endpoint
* [ ] Network & volume tanımları
* [ ] Basit deploy dokümantasyonu

### 📦 Çıktılar

* `docker-compose up` ile çalışan sistem
* Self-hosted production setup

---

# 🤖 FAZ 8 – GELİŞMİŞ ZEKÂ (Sprint 8)

### 🎯 Amaç

Rakiplerden ayrışma

### 🔹 Task List

* [ ] Service behavior profiling
* [ ] Adaptive threshold
* [ ] Versiyon bazlı performans karşılaştırma
* [ ] Otomatik "rollback uyarısı" önerisi

### 📦 Çıktılar

* Akıllı uyarılar
* Pro seviyesinde özellikler

---

## 🏁 SON DURUM

Bu sprint planı sonunda:

* Gerçek müşteriye satılabilir
* SaaS veya self-hosted
* CV + ticari ürün

---

**Bir sonraki adım:**

* İstersen her sprint için **ayrı ayrı teknik detay + Go kod iskeleti** çıkarırım.
