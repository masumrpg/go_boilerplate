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
│   │   │   └── migration.go       # Database migration
│   │   │
│   │   ├── middleware/
│   │   │   ├── auth.go            # JWT middleware
│   │   │   ├── logger.go          # Logging middleware
│   │   │   ├── cors.go            # CORS middleware
│   │   │   └── validator.go       # Request validator middleware
│   │   │
│   │   └── utils/
│   │       ├── jwt.go             # JWT token utilities
│   │       ├── hash.go            # Password hashing (bcrypt)
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
- Initialize shared components (config, database, logger)
- Register all module routes
- Start Fiber server

### **`internal/shared/`** - Shared Components

#### `config/`
- Load configuration dari environment variables
- Menggunakan Viper
- Config struct untuk type-safe access

#### `database/`
- Database connection pooling (GORM + PostgreSQL)
- Auto migration
- Connection management

#### `middleware/`
- **auth.go**: JWT validation middleware
- **logger.go**: HTTP request/response logging
- **cors.go**: CORS configuration
- **validator.go**: Request body validation

#### `utils/`
- **jwt.go**: Generate & validate JWT tokens
- **hash.go**: Password hashing dengan bcrypt
- **validator.go**: Struct validation helper
- **response.go**: Standardized JSON response
- **logger.go**: Logrus configuration

---

### **`internal/modules/`** - Feature Modules

#### Standard Module Structure
Setiap module mengikuti pattern yang sama:

```
module-name/
├── model.go           # Domain entity (GORM model)
├── repository.go      # Data access layer (database operations)
├── service.go         # Business logic layer
├── handler.go         # HTTP request handlers (controller)
├── routes.go          # Route registration & middleware setup
└── dto/
    ├── request.go     # Input validation DTOs
    └── response.go    # Output DTOs
```

#### Layer Responsibilities

**`model.go`**
- Define database entity/schema
- GORM struct tags
- Table relationships
- Hooks (BeforeCreate, AfterUpdate, etc)

**`repository.go`**
- Interface definition
- CRUD operations
- Database queries
- No business logic

**`service.go`**
- Interface definition
- Business logic implementation
- Orchestrate multiple repositories
- Call external services
- Data transformation

**`handler.go`**
- Parse HTTP request
- Validate input (call validator)
- Call service methods
- Format & return HTTP response
- Error handling

**`routes.go`**
- Register routes untuk module
- Apply middleware (auth, validator, etc)
- Group related endpoints
- Dependency injection dari main.go

**`dto/request.go`**
- Input validation structs
- Validation tags (required, email, min, max, etc)
- Request body parsing

**`dto/response.go`**
- Output format structs
- Hide sensitive fields (password, dll)
- Consistent response structure

---

## 🔄 Request Flow Diagram

```mermaid
graph TD
    A[HTTP Request] --> B[Fiber Router]
    B --> C[Global Middleware]
    C --> C1[Logger Middleware]
    C1 --> C2[CORS Middleware]
    C2 --> D[Module Routes]

    D --> E[Route Middleware]
    E --> E1[JWT Auth Middleware]
    E1 --> E2[Validator Middleware]

    E2 --> F[Handler Layer]
    F --> F1[Parse Request]
    F1 --> F2[Validate DTO]
    F2 --> F3[Call Service]

    F3 --> G[Service Layer]
    G --> G1[Business Logic]
    G1 --> G2[Call Repository]
    G1 --> G3[Call External Service]

    G2 --> H[Repository Layer]
    H --> H1[Database Query]
    H1 --> I[(PostgreSQL)]

    I --> H2[Return Model]
    H2 --> G4[Transform to DTO]

    G3 --> J[External Services]
    J --> J1[Email Service]
    J --> J2[OAuth Service]

    G4 --> F4[Format Response]
    F4 --> K[HTTP Response JSON]

    style A fill:#e1f5ff
    style K fill:#e1f5ff
    style F fill:#fff4e1
    style G fill:#f0e1ff
    style H fill:#e1ffe1
    style C fill:#ffe1e1
    style E fill:#ffe1e1
```

---

## 🔀 Module Internal Flow

```mermaid
graph LR
    A[routes.go] -->|Register Routes| B[handler.go]
    B -->|Call Business Logic| C[service.go]
    C -->|Query Database| D[repository.go]
    D -->|GORM Operations| E[(Database)]

    B -->|Parse & Validate| F[dto/request.go]
    C -->|Transform Data| G[dto/response.go]
    G -->|Return to| B

    D -->|Use Entity| H[model.go]

    style A fill:#e1f5ff
    style B fill:#fff4e1
    style C fill:#f0e1ff
    style D fill:#e1ffe1
    style E fill:#ffcccc
```

---

## 🔌 Dependency Injection Flow

```mermaid
graph TD
    A[main.go] --> B[Initialize Config]
    A --> C[Initialize Database]
    A --> D[Initialize Logger]

    B --> E[Create Fiber App]
    C --> E
    D --> E

    E --> F[Register Module Routes]

    F --> G[auth.RegisterRoutes]
    F --> H[user.RegisterRoutes]
    F --> I[email.RegisterRoutes]
    F --> J[oauth.RegisterRoutes]

    G --> G1[Create Repository]
    G1 --> G2[Create Service]
    G2 --> G3[Create Handler]
    G3 --> G4[Register Routes]

    H --> H1[Create Repository]
    H1 --> H2[Create Service]
    H2 --> H3[Create Handler]
    H3 --> H4[Register Routes]

    style A fill:#ff9999
    style E fill:#99ccff
    style G4 fill:#99ff99
    style H4 fill:#99ff99
```

---

## 📊 Module Interaction Flow

```mermaid
graph TD
    A[Client Request] --> B{Route Type}

    B -->|/auth/login| C[Auth Module]
    B -->|/auth/register| C

    B -->|/users/*| D[User Module]

    B -->|/oauth/*| E[OAuth Module]

    C -->|Generate JWT| G[shared/utils/jwt.go]
    C -->|Hash Password| H[shared/utils/hash.go]
    C -->|Send Email| I[Email Module]

    D -->|Validate Token| G
    E -->|Validate Token| G
    E -->|Send Notification| I

    style A fill:#e1f5ff
    style G fill:#ffe1e1
    style H fill:#ffe1e1
    style I fill:#f0e1ff
```

---

## 🎯 Layer Communication Pattern

```mermaid
graph TB
    subgraph "Module A (User)"
        A1[Handler A] --> A2[Service A] --> A3[Repository A]
    end

    subgraph "Module B (Email)"
        B1[Service B - No Repository]
    end

    subgraph "Shared Components"
        S1[Middleware]
        S2[Utils]
        S3[Config]
        S4[(Database)]
    end

    A1 -.Use.-> S1

    A2 -.Use.-> S2

    A3 --> S4

    A2 -->|Send Email| B1

    style S1 fill:#ffe1e1
    style S2 fill:#ffe1e1
    style S3 fill:#ffe1e1
    style S4 fill:#ffcccc
```

---

## 📋 Naming Convention

### Files
- **model.go** - Singular noun (User, Product, Order)
- **repository.go** - Data access methods
- **service.go** - Business logic methods
- **handler.go** - HTTP handlers
- **routes.go** - Route registration

### Interfaces & Structs
- **Interface**: `UserRepository`, `UserService`
- **Implementation**: `userRepository`, `userService` (private)
- **Constructor**: `NewRepository()`, `NewService()`, `NewHandler()`

### Methods
- **Repository**: `Create`, `FindByID`, `Update`, `Delete`, `FindAll`
- **Service**: Business-specific names (`GetUserProfile`, `CreateOrder`, `ProcessPayment`)
- **Handler**: HTTP verb-like names (`GetUser`, `CreateUser`, `UpdateUser`, `DeleteUser`)

---

## 🚀 Module Development Workflow

```mermaid
graph LR
    A[1. Create Model] --> B[2. Create Repository]
    B --> C[3. Create Service]
    C --> D[4. Create DTOs]
    D --> E[5. Create Handler]
    E --> F[6. Create Routes]
    F --> G[7. Register in main.go]

    style A fill:#e1ffe1
    style G fill:#ff9999
```

---

## ✨ Keuntungan Pattern Ini

### 1. **Clear Separation of Concerns**
- Shared components terpisah dari module-specific code
- Setiap layer punya tanggung jawab jelas

### 2. **Scalability**
- Tambah module baru tanpa ganggu existing code
- Module independen satu sama lain

### 3. **Maintainability**
- Mudah cari kode (semua tentang User ada di `modules/user/`)
- Perubahan di satu module tidak affect module lain

### 4. **Testability**
- Mudah mock dependencies per layer
- Test isolated per module

### 5. **Team Collaboration**
- Developer bisa kerja di module berbeda
- Minimal merge conflicts

### 6. **Reusability**
- Shared utils bisa dipakai semua module
- No code duplication

---

## 📦 Example Module List

```
internal/
├── shared/              # Komponen global
│   ├── config/
│   ├── database/
│   ├── middleware/
│   └── utils/
│
└── modules/             # Feature modules
    ├── auth/            # Authentication & authorization
    ├── user/            # User management
    ├── email/           # Email notifications
    └── oauth/           # Social login
```

---

## 🔧 Quick Start: Add New Module

```bash
# Create module structure
mkdir -p internal/modules/payment/dto
cd internal/modules/payment

# Create files
touch model.go repository.go service.go handler.go routes.go
touch dto/request.go dto/response.go
```

Implement pattern yang sama seperti module lain, lalu register di `main.go`!

---

## 📚 Technology Stack

- **Framework**: Fiber v2
- **ORM**: GORM + PostgreSQL Driver
- **Validation**: go-playground/validator/v10
- **JWT**: golang-jwt/jwt/v5 + gofiber/contrib/jwt
- **JSON**: bytedance/sonic (fast serialization)
- **Config**: spf13/viper
- **Logger**: sirupsen/logrus
- **Email**: gopkg.in/gomail.v2
- **OAuth**: golang.org/x/oauth2
- **Docs**: swaggo/swag + gofiber/swagger
- **Testing**: stretchr/testify