# 🔗 Smolink

**Smolink** is a production-grade URL shortening service built in **Go**, supporting:

- ✅ Custom aliases & automatic short codes  
- 📊 Click tracking & analytics  
- 🛎️ Webhook delivery with retries  
- ⚡ Redis caching  
- 🧼 Clean architecture & layered design  
- 🧪 Full test suite with unit and integration testing  

---

## 🚀 Features

- **Shorten URLs** with or without custom aliases
- **Resolve URLs** and automatically track:
  - Click count
  - IP address
  - User agent
  - Timestamps
- **Webhook support** with retry mechanism
- **In-memory + Redis cache** for speed
- **PostgreSQL** as primary data store
- Graceful startup & shutdown
- Structured logging middleware
- Modular & testable codebase

---

## 📁 Project Structure

```
smolink/
├── cmd/                  # Main application entry point
├── internal/
│   ├── app/              # App wiring (router, DI)
│   ├── config/           # Env & config loader
│   ├── controller/       # HTTP route handlers
│   ├── service/          # Core business logic
│   ├── repository/       # PostgreSQL & Redis access
│   ├── model/            # Database and request models
│   └── migration/        # Database schema setup
├── pkg/
│   ├── database/         # DB and Redis initialization
│   └── logger/           # Logging middleware
├── test/                 # Integration testing with Dockertest
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Requirements

- Go 1.21+
- Docker (for integration tests)
- PostgreSQL
- Redis

---

## 🧪 Run Tests

Smolink supports both unit and integration tests.

To run all tests:

```bash
go test -v -race ./...
```

Integration tests use Dockerized PostgreSQL and Redis via [dockertest](https://github.com/ory/dockertest).

---

## 🏗️ Run Locally

### 1. Environment Setup

Create a `.env` file or export the following environment variables:

```env
ENV=development
SERVER_PORT=:8080
POSTGRES_DSN=postgres://user:pass@localhost:5432/smolink?sslmode=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 2. Start PostgreSQL & Redis

```bash
docker run --rm -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=smolink postgres:15-alpine

docker run --rm -p 6379:6379 redis:7-alpine
```

### 3. Run the app

```bash
go run cmd/main.go
```

---

## 📫 API Endpoints

| Method | Endpoint       | Description             |
|--------|----------------|-------------------------|
| POST   | `/shorten`     | Shorten a URL           |
| GET    | `/:short_code` | Redirect to full URL    |

### Sample Request (POST `/shorten`)

```json
{
  "url": "https://example.com",
  "customAlias": "my-custom-code"
}
```

---

## 🛠 Developer Utilities

Format and vet your code:

```bash
go fmt ./...
go vet ./...
```

---

## 📦 CI & Testing

This project includes GitHub Actions to automatically run tests on each push and PR.

---
