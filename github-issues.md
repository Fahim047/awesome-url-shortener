# 📌 GitHub Issues for Awesome URL Shortener

## Issue 1: Project Initialization

**Description:**\
Set up the Go project with basic scaffolding.

**Tasks:**\

- \[ \] Initialize Git repo\
- \[ \] Initialize Go module
  (`go mod init github.com/<user>/url-shortener`)\
- \[ \] Add `.gitignore` (ignore `/vendor`, binaries, IDE files)\
- \[ \] Create `main.go` with a basic `net/http` server\
- \[ \] Add `/healthz` endpoint

---

## Issue 2: Database Setup (PostgreSQL)

**Description:**\
Implement database schema and connection.

**Tasks:**\

- \[ \] Define DB connection config via environment variables\
- \[ \] Connect to Postgres (`pgx` or `pq`)\
- \[ \] Write migration to create `url_mappings` table\
- \[ \] Add index on `short_key`\
- \[ \] Create `pkg/db` package with functions:\
- \[ \] `CreateMapping`\
- \[ \] `GetMapping`\
- \[ \] `IncrementClickCount`

---

## Issue 3: Cache Layer (Redis)

**Description:**\
Add caching to reduce DB lookups and track clicks.

**Tasks:**\

- \[ \] Define Redis config via environment variables\
- \[ \] Connect using `go-redis` client\
- \[ \] Implement:\
- \[ \] `CacheSet(shortKey, longURL, ttl)`\
- \[ \] `CacheGet(shortKey)`\
- \[ \] `CacheIncrClicks(shortKey)`

---

## Issue 4: Short Key Generation

**Description:**\
Create logic to generate short keys.

**Tasks:**\

- \[ \] Implement Base62 encoder (id → short key)\
- \[ \] Add random key generator fallback with collision check\
- \[ \] Validate custom alias availability

---

## Issue 5: API Endpoint - Shorten URL

**Description:**\
Create a new short URL.

**Tasks:**\

- \[ \] Implement `POST /api/v1/shorten`\
- \[ \] Parse JSON body: `long_url`, `custom_alias`, `expire_at`\
- \[ \] Validate URL format\
- \[ \] Insert mapping into DB\
- \[ \] Cache entry in Redis\
- \[ \] Return `{ "short_url": "http://domain/<shortKey>" }`

---

## Issue 6: API Endpoint - Redirect

**Description:**\
Redirect users from short URL → long URL.

**Tasks:**\

- \[ \] Implement `GET /:short_key`\
- \[ \] Lookup in Redis first\
- \[ \] Fallback to DB if cache miss\
- \[ \] Handle expired URLs\
- \[ \] Increment click count in Redis\
- \[ \] Return `302 Found`

---

## Issue 7: API Endpoint - Analytics (Optional)

**Description:**\
Provide analytics for a given short key.

**Tasks:**\

- \[ \] Implement `GET /api/v1/analytics/:short_key`\
- \[ \] Query Postgres for `click_count`, `created_at`, `expire_at`\
- \[ \] Merge Redis counters if needed\
- \[ \] Return JSON
  `{ short_key, long_url, click_count, created_at, expire_at }`

---

## Issue 8: API Endpoint - Top URLs (Optional)

**Description:**\
Return most clicked URLs.

**Tasks:**\

- \[ \] Implement `GET /api/v1/top`\
- \[ \] Query DB for top 10 `short_key` ordered by `click_count`\
- \[ \] Return JSON list

---

## Issue 9: Click Counting Sync Job

**Description:**\
Keep DB click counts in sync with Redis.

**Tasks:**\

- \[ \] Increment Redis counter on redirect\
- \[ \] Background goroutine every 5 minutes:\
- \[ \] Read counters from Redis\
- \[ \] Update Postgres with increments\
- \[ \] Reset counters in Redis

---

## Issue 10: Docker & Local Dev

**Description:**\
Containerize app and setup local development.

**Tasks:**\

- \[ \] Write Dockerfile (multi-stage build)\
- \[ \] Add `.dockerignore`\
- \[ \] Create `docker-compose.yml` with services:\
- \[ \] Go app\
- \[ \] Postgres\
- \[ \] Redis

---

## Issue 11: CI/CD Pipeline

**Description:**\
Setup GitHub Actions for linting, testing, and building.

**Tasks:**\

- \[ \] Add workflow `.github/workflows/ci.yml`\
- \[ \] Run `golangci-lint run ./...`\
- \[ \] Run `go test ./...`\
- \[ \] Build Docker image\
- \[ \] Push to registry (optional)

---

## Issue 12: Testing

**Description:**\
Add unit, integration, and load tests.

**Tasks:**\

- \[ \] Unit tests:\
- \[ \] Short key generation\
- \[ \] URL validation\
- \[ \] DB mocks\
- \[ \] Redis mocks\
- \[ \] Integration tests (with Dockerized Postgres + Redis)\
- \[ \] Load testing (optional, with `hey` or `wrk`)

---

## Issue 13: Logging & Monitoring

**Description:**\
Add observability to the service.

**Tasks:**\

- \[ \] Structured JSON logging (`zerolog` or `logrus`)\
- \[ \] Add `/metrics` endpoint for Prometheus\
- \[ \] Track request counts, latency, cache hit/miss

---

## Issue 14: Documentation

**Description:**\
Write documentation for developers.

**Tasks:**\

- \[ \] Write `README.md` (overview, setup, API docs)\
- \[ \] Add environment variable reference\
- \[ \] Add API endpoint examples\
- \[ \] (Optional) Add OpenAPI/Swagger spec
