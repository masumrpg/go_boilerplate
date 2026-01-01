# 🛡️ Go Boilerplate - Enterprise-Ready Modular API

Sebuah boilerplate REST API yang kokoh, modular (feature-based), dan siap produksi menggunakan ekosistem Go modern. Dirancang untuk skalabilitas, keamanan, dan kemudahan pengembangan.

---

## 🚀 Fitur Utama

- **Modular Architecture**: Struktur folder berbasis fitur (Domain-Driven Design friendly).
- **Advanced Auth System**:
  - JWT Authentication (Access & Refresh Tokens).
  - RBAC (Role-Based Access Control) dengan permission granular.
  - Multi-factor Authentication (2FA) & Verifikasi Email.
  - Session Management & Device Tracking.
- **OAuth2 Integration**: Login via Google & GitHub.
- **Robust Persistence**: GORM dengan dukungan PostgreSQL.
- **Caching & OTP**: Redis untuk validasi OTP yang cepat dan aman.
- **Embedded Email Templates**: Template HTML yang dinamis dengan `//go:embed`.
- **Automatic Swagger**: Dokumentasi API interaktif yang selalu sinkron.
- **Database Migrations**: Manajemen skema versi menggunakan `golang-migrate`.
- **Docker Ready**: Deployment instan dengan Docker Compose.

---

## 🏗️ Arsitektur Sistem

Aplikasi ini menggunakan pola **Modular Layered Architecture**. Setiap modul merangkum logikanya sendiri sementara tetap berbagi komponen universal di folder `shared`.

```mermaid
graph TD
    subgraph Client_Layer
        User([User Client])
    end

    subgraph API_Layer
        Fiber[Fiber v2 Framework]
        Handler[Handlers/Controllers]
        M_Auth[JWT Middleware]
        M_RBAC[RBAC Middleware]
        M_Log[Logger Middleware]
    end

    subgraph Logic_Layer
        Service[Service / Business Logic]
        DTO[DTO - Data Transfer Object]
    end

    subgraph Data_Layer
        Repo[Repository / Data Access]
        GORM[GORM ORM]
    end

    subgraph External_Resources
        PG[(Postgres Database)]
        RD[(Redis Cache/OTP)]
        SMTP[Gomail SMTP]
    end

    User -- HTTP/JSON --> Fiber
    Fiber --> M_Log
    M_Log --> M_Auth
    M_Auth --> M_RBAC
    M_RBAC --> Handler
    Handler --> Service
    Service -- Validation --> DTO
    Service --> Repo
    Repo --> GORM
    GORM --> PG
    Service --> RD
    Service --> SMTP
```

---

## 🔐 Alur Autentikasi & Keamanan

Berikut adalah alur pendaftaran hingga login dengan fitur keamanan berlapis:

```mermaid
sequenceDiagram
    participant User
    participant API as API Server
    participant Redis
    participant DB as PostgreSQL
    participant Mail as Email Service

    Note over User, Mail: Registrasi Account
    User->>API: POST /register (Data User)
    API->>DB: Save User (Status: Unverified)
    API->>Redis: Set Activation CODE (TTL 10m)
    API->>Mail: Send Verification Email
    User->>API: POST /verify-email (CODE)
    API->>Redis: Validate CODE
    API->>DB: Update User (is_verified: true)

    Note over User, Mail: Login Flow with 2FA
    User->>API: POST /login (Credentials)
    API->>DB: Verify Password
    alt 2FA is Enabled
        API->>Redis: Set 2FA CODE (TTL 5m)
        API->>Mail: Send 2FA OTP Email
        User->>API: POST /verify-2fa (CODE)
        API->>Redis: Validate CODE
    end
    API->>DB: Create Session (Metadata: IP, Device)
    API-->>User: Return JWT (Access & Refresh)
```

---

## 🛠️ Technology Stack & Penggunaan

| Komponen | Library | Alasan & Penggunaan |
| :--- | :--- | :--- |
| **Web Framework** | [Fiber v2](https://gofiber.io/) | Framework Go tercepat (berbasis fasthttp) dengan performa tinggi & middleware lengkap. |
| **Database ORM** | [GORM](https://gorm.io/) | ORM paling populer di Go. Digunakan untuk query, relasi, dan auto-migration. |
| **Database** | [PostgreSQL](https://www.postgresql.org/) | RDBMS powerfull untuk konsistensi data dan integritas. |
| **Cache & OTP** | [Redis](https://redis.io/) | Digunakan untuk menyimpan kode OTP activation/2FA dengan TTL dan session tracking. |
| **Auth** | [JWT-Go (v5)](https://github.com/golang-jwt/jwt) | Implementasi token keamanan berbasis standar industri. |
| **Configuration** | [Viper](https://github.com/spf13/viper) | Membaca konfigurasi dari `.env`, env vars, atau config file secara dinamis. |
| **Validation** | [Validator v10](https://github.com/go-playground/validator) | Validasi request DTO (email, required, min-max) menggunakan struct tags. |
| **Logging** | [Logrus](https://github.com/sirupsen/logrus) | Structured logging dengan level (info, warn, error) dan format JSON/Text. |
| **Email** | [Gomail](https://github.com/go-gomail/gomail) | Mengelola pengiriman email SMTP untuk notifikasi dan OTP. |
| **Migrations** | [Golang-Migrate](https://github.com/golang-migrate/migrate) | Versioning database schema secara eksplisit dan aman. |
| **API Docs** | [Swaggo](https://github.com/swaggo/swag) | Meng-generate OpenAPI 2.0 dokumentasi langsung dari code comments. |

---

## 📁 Struktur Folder Modular

```text
internal/
├── shared/                # 🛠️ GLOBAL COMPONENTS
│   ├── config/            # Viper setup
│   ├── database/          # Connections (GORM, Redis) & Migrations
│   ├── middleware/        # Auth, RBAC, Logger, CORS, Validator
│   └── utils/             # JWT, Hash, Response, Logger Helpers
│
└── modules/               # 🔥 DOMAIN MODULES
    ├── auth/              # Logic Login, Registration, 2FA
    ├── user/              # Management User & Profile
    ├── role/              # RBAC: Create Role & Permissions
    ├── oauth/             # Google & GitHub Login
    └── email/             # SMTP Service & HTML Templates
```

---

## 🛠️ Cara Menjalankan

### Menggunakan Docker (Rekomendasi)
```bash
# 1. Clone repository
# 2. Setup .env (copy dari .env.example)
docker-compose up -d --build
```

### Manual (Development)
```bash
# 1. Install dependencies
go mod download

# 2. Run migrations
make migrate-up

# 3. Generate Swagger (jika ada perubahan handler)
make swagger

# 4. Start Server
go run cmd/api/main.go
```

---

## 📚 API Dokumentasi

Setelah server berjalan, dokumentasi lengkap tersedia di:
- **Swagger UI**: `http://localhost:3000/swagger/`
- **Postman**: Import file `Go_Boilerplate_API.postman_collection.json` di root folder.

---

## 🔐 Keamanan & Fitur Tambahan

- **Device Fingerprinting**: Setiap session mencatat `IP Address`, `User Agent`, dan `Device ID`.
- **RBAC Granular**: Akses dikontrol hingga tingkat resource (contoh: `roles:create`, `users:update`).
- **Clean Shutdown**: Menangani signal Linux (`SIGINT`, `SIGTERM`) untuk menutup koneksi database secara aman saat server mati.
- **Embedded Templates**: Template email berada di dalam binary, memudahkan distribusi tanpa perlu mengcopy file manual.

---
© 2026 Go Boilerplate Squad.
