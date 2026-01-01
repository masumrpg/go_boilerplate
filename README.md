# Struktur Project Golang Modular (Feature-Based)

## 📁 Struktur Folder

```
project-root/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point aplikasi
│
├── internal/
│   ├── shared/                    # 🔧 SHARED COMPONENTS
│   │   ├── config/
│   │   │   └── config.go          # Konfigurasi (Viper)
│   │   │
│   │   ├── database/
│   │   │   ├── connection.go      # Database connection (GORM + PostgreSQL)
│   │   │   ├── redis.go           # Redis connection
│   │   │   └── migration.go       # Database migration & table rename
│   │   │
│   │   ├── middleware/
│   │   │   ├── auth.go            # JWT middleware + RBAC
│   │   │   ├── logger.go          # Logging middleware
│   │   │   ├── cors.go            # CORS middleware
│   │   │   └── validator.go       # Request validator middleware
│   │   │
│   │   └── utils/
│   │       ├── jwt.go             # JWT token utilities
│   │       ├── hash.go            # Password hashing (bcrypt)
│   │       ├── random.go          # Random string helper (OTP)
│   │       ├── validator.go       # Struct validation helper
│   │       ├── response.go        # Standard API response format
│   │       └── logger.go          # Logger setup & helper
│   │
│   └── modules/                   # 🔥 FEATURE MODULES
│       │
│       ├── auth/                  # AUTH MODULE
│       │   ├── model.go           # Auth-related models (jika ada)
│       │   ├── repository.go      # Auth data access
│       │   ├── service.go         # Auth business logic
│       │   ├── handler.go         # Auth HTTP handlers
│       │   ├── routes.go          # Auth route registration
│       │   └── dto/
│       │       ├── request.go     # Login, Register, Refresh DTOs
│       │       └── response.go    # Token response DTOs
│       │
│       ├── user/                  # USER MODULE
│       │   ├── model.go           # User entity/model
│       │   ├── repository.go      # User repository (CRUD)
│       │   ├── service.go         # User business logic
│       │   ├── handler.go         # User HTTP handlers
│       │   ├── routes.go          # User route registration
│       │   └── dto/
│       │       ├── request.go     # Create, Update user DTOs
│       │       └── response.go    # User response DTOs
│       │
│       ├── role/                  # ROLE MODULE (RBAC)
│       │   ├── model.go           # Role entity/model
│       │   ├── repository.go      # Role repository
│       │   ├── service.go         # Role business logic + seeding
│       │   ├── handler.go         # Role HTTP handlers
│       │   ├── routes.go          # Role route registration
│       │   └── dto/
│       │       ├── request.go     # Create, Update role DTOs
│       │       └── response.go    # Role response DTOs
│       │
│       ├── email/                 # EMAIL MODULE
│       │   ├── service.go         # Email service (gomail)
│       │   ├── template.go        # Email HTML templates
│       │   └── dto/
│       │       └── request.go     # Email send request DTO
│       │
│       └── oauth/                 # OAUTH MODULE
│           ├── service.go         # OAuth2 service (Google, GitHub)
│           ├── handler.go         # OAuth callback handlers
│           ├── routes.go          # OAuth routes
│           └── dto/
│               └── response.go    # OAuth user info response
│
├── docs/
│   ├── docs.go                    # Swagger generated files
│   ├── swagger.json
│   └── swagger.yaml
│
├── pkg/
│   └── ...                        # Public packages (optional)
│
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 🏗️ Pattern & Responsibility

### **`cmd/api/main.go`**
- Initialize shared components (config, database, logger, redis)
- Register all module routes
- Start Fiber server

### **`internal/shared/`** - Shared Components

#### `config/`
- Load configuration dari environment variables (Viper)
- Config struct untuk type-safe access

#### `database/`
- Database connection pooling (GORM + PostgreSQL)
- Redis client configuration
- Migration management

#### `middleware/`
- **auth.go**: JWT validation middleware & RBAC
- **logger.go**: HTTP request/response logging
- **cors.go**: CORS configuration
- **validator.go**: Request body validation

#### `utils/`
- **jwt.go**: Generate & validate JWT tokens
- **hash.go**: Password hashing dengan bcrypt
- **random.go**: Random string generator untuk OTP
- **validator.go**: Struct validation helper

---

### 🔒 Security Features (Optional)

Aplikasi ini mendukung fitur keamanan tambahan yang bisa diaktifkan melalui `.env`:

#### 1. Account Activation (Email Verification)
- **Env**: `EMAIL_VERIFICATION_ENABLED=true`
- **Deskripsi**: User baru harus memverifikasi email dengan kode OTP 6-digit sebelum bisa login.
- **Penyimpanan**: Kode OTP disimpan di Redis (TTL 10 menit).

#### 2. Two-Factor Authentication (2FA)
- **Env**: `TWO_FACTOR_ENABLED=true`
- **Deskripsi**: Setelah memasukkan password, user harus memasukkan kode OTP yang dikirim ke email.
- **Penyimpanan**: Kode OTP disimpan di Redis (TTL 5 menit).

---

## 💾 Database Migrations

Menggunakan `golang-migrate` untuk manajemen skema database yang versi-able.

- **Run Migrations**: `make migrate-up` atau `go run cmd/migrate/main.go -up`
- **Rollback**: `make migrate-down`
- **Create New**: `make migrate-create`

---

## 🚀 Docker Support

Aplikasi sudah mendukung containerization:
- **Run**: `docker-compose up -d --build`
- **Services**: App, PostgreSQL, Redis, Migrate.

---

## 🏗️ Architecture Layers

### **`internal/modules/`** - Feature Modules

Setiap module mengikuti pattern yang sama:

**`model.go`**
- Define database entity/schema, GORM struct tags, relationships.

**`repository.go`**
- Interface & implementation untuk data access (queries only).

**`service.go`**
- Business logic implementation, orchestrate repositories.

**`handler.go`**
- Parse HTTP requests, call service methods, return responses.

**`routes.go`**
- Register routes, apply middleware, dependency injection.

---

## 📚 Technology Stack

- **Framework**: Fiber v2
- **ORM**: GORM + PostgreSQL
- **Caching/OTP**: Redis
- **Validation**: go-playground/validator/v10
- **JWT**: golang-jwt/jwt/v5
- **Logger**: sirupsen/logrus
- **Email**: gopkg.in/gomail.v2
- **OAuth**: golang.org/x/oauth2

---

## 🔐 RBAC System

API ini menggunakan sistem RBAC (Role-Based Access Control):
- **3 Default Role**: SuperAdmin, Admin, User.
- **Granular Permissions**: Format `resource.action` (contoh: `users.create`).
- **SuperAdmin Account**: Otomatis dibuat saat startup berdasarkan `.env`.

---

## 🎛️ Feature Flags

Fitur opsional via `.env`:
- `OAUTH_GOOGLE_ENABLED`: Aktifkan Google OAuth.
- `OAUTH_GITHUB_ENABLED`: Aktifkan GitHub OAuth.
- `EMAIL_ENABLED`: Aktifkan pengiriman email.
- `EMAIL_VERIFICATION_ENABLED`: Aktifkan verifikasi email user baru.
- `TWO_FACTOR_ENABLED`: Aktifkan 2FA login.
