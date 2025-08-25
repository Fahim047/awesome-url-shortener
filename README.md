# Awesome URL Shortener (Go + PostgreSQL + Redis)

A simple URL shortener service built with **Go (net/http)**, **PostgreSQL**, and **Redis**.  
It supports database-backed URL mappings, caching, and a health check endpoint.

---

## 🚀 Features (Implemented)

- REST API server with **Go (net/http)**.
- **PostgreSQL** schema for storing URL mappings.
- **Redis** cache for fast lookups.
- Graceful shutdown handling (`SIGINT`, `SIGTERM`).

---

## 🛠️ Tech Stack

- **Go 1.24+**
- **PostgreSQL 17 (alpine)**
- **Redis 7 (alpine)**
- **Docker & docker-compose**

---

## ⚙️ Setup & Run

### 1. Clone the repo

```bash
git clone https://github.com/Fahim047/awesome-url-shortener.git
cd awesome-url-shortener
```

### 2. Start dependencies (Postgres + Redis)

```bash
docker compose up -d
```

- Postgres → `localhost:5432` (user: `postgres`, pass: `postgres`, db: `url_shortener`)
- Redis → `localhost:6379`

The database schema is auto-loaded from `schema.sql`.

### 3. Run the Go app (using make)

```bash
make run
```

---

## 📂 Project Structure

```
awesome-url-shortener/
│── cmd/
│   └── server/       # main.go (entrypoint)
│── pkg/              # (db, cache, api)
│── docker-compose.yaml
│── schema.sql
│── go.mod / go.sum
```

---

## ✅ Endpoints

- `POST /api/v1/shorten` → Create short URL.
- `GET /:short_key` → Redirect to long URL.
- `GET /api/v1/analytics/:short_key` → Fetch analytics.
- `GET /api/v1/top` → Fetch top 10 links.

---
