# Chirpy 🐦

A Twitter-like microblogging REST API built in Go, developed as part of the [Boot.dev](https://boot.dev) *Learn HTTP Servers in Go* course.

## Overview

Chirpy lets users register, post short messages called **chirps** (max 140 characters), and manage their accounts — all secured with JWT-based authentication and Argon2id password hashing. It also supports a premium membership tier called **Chirpy Red**, upgradeable via a webhook from the Polka payment provider.

## Tech Stack

- **Language:** Go (standard library `net/http`)
- **Database:** PostgreSQL
- **ORM / Query Layer:** [sqlc](https://sqlc.dev) — type-safe SQL queries
- **Migrations:** [Goose](https://github.com/pressly/goose)
- **Auth:** JWT (`golang-jwt/jwt`) + Argon2id (`alexedwards/argon2id`)
- **UUID:** `google/uuid`

## Prerequisites

- Go 1.22+
- PostgreSQL
- [Goose](https://github.com/pressly/goose) (for running migrations)
- [sqlc](https://sqlc.dev) (only if regenerating query code)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Shinzi04/boot-go-http-server.git
cd boot-go-http-server
```

### 2. Set up the database

Create a PostgreSQL database and run the migrations:

```bash
goose -dir sql/schema postgres "<your-db-url>" up
```

### 3. Configure environment variables

Create a `.env` file in the project root:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=your_jwt_secret_here
POLKA_KEY=your_polka_webhook_key_here
```

| Variable    | Description                                              |
|-------------|----------------------------------------------------------|
| `DB_URL`    | PostgreSQL connection string                             |
| `PLATFORM`  | Deployment environment (`dev` enables the reset endpoint)|
| `JWT_SECRET`| Secret key used to sign JWT access tokens                |
| `POLKA_KEY` | API key expected from the Polka webhook                  |

### 4. Run the server

```bash
go run .
```

The server starts on **`http://localhost:8080`**.

---

## API Reference

### Health

| Method | Endpoint        | Description        | Auth     |
|--------|-----------------|--------------------|----------|
| GET    | `/api/healthz`  | Readiness check    | None     |

### Users

| Method | Endpoint      | Description                     | Auth         |
|--------|---------------|---------------------------------|--------------|
| POST   | `/api/users`  | Register a new user             | None         |
| PUT    | `/api/users`  | Update email & password         | Bearer JWT   |

**Register** — `POST /api/users`
```json
{ "email": "user@example.com", "password": "secret" }
```

### Authentication

| Method | Endpoint       | Description                              | Auth          |
|--------|----------------|------------------------------------------|---------------|
| POST   | `/api/login`   | Login and receive tokens                 | None          |
| POST   | `/api/refresh` | Exchange refresh token for access token  | Bearer (refresh token) |
| POST   | `/api/revoke`  | Revoke a refresh token                   | Bearer (refresh token) |

**Login** — `POST /api/login`
```json
{ "email": "user@example.com", "password": "secret" }
```
Returns a JWT access token (expires in **1 hour**) and a refresh token (expires in **60 days**).

### Chirps

| Method | Endpoint                   | Description                        | Auth       |
|--------|----------------------------|------------------------------------|------------|
| POST   | `/api/chirps`              | Create a new chirp                 | Bearer JWT |
| GET    | `/api/chirps`              | List all chirps                    | None       |
| GET    | `/api/chirps/{chirpID}`    | Get a single chirp by ID           | None       |
| DELETE | `/api/chirps/{chirpID}`    | Delete a chirp (owner only)        | Bearer JWT |

**Create Chirp** — `POST /api/chirps`
```json
{ "body": "Hello world!" }
```
- Maximum 140 characters.
- Profane words (`kerfuffle`, `sharbert`, `fornax`) are automatically replaced with `****`.

**List Chirps** — `GET /api/chirps`

Supports optional query parameters:
- `author_id` — filter chirps by user UUID
- `sort` — `asc` (default) or `desc` by creation date

### Webhooks

| Method | Endpoint                  | Description                         | Auth   |
|--------|---------------------------|-------------------------------------|--------|
| POST   | `/api/polka/webhooks`     | Upgrade user to Chirpy Red          | ApiKey |

Expects `Authorization: ApiKey <POLKA_KEY>` and a payload with `event: "user.upgraded"`.

### Admin

| Method | Endpoint          | Description                          | Auth |
|--------|-------------------|--------------------------------------|------|
| GET    | `/admin/metrics`  | View file server hit count           | None |
| POST   | `/admin/reset`    | Reset hit counter & users (dev only) | None |

### Static Files

Files in the project root are served under `/app/`.

---

## Project Structure

```
.
├── main.go                        # Server setup, routing
├── handler_chirps_create.go       # POST /api/chirps
├── handler_chirps_get.go          # GET /api/chirps[/{id}]
├── handler_chirps_delete.go       # DELETE /api/chirps/{id}
├── handler_users_create.go        # POST /api/users
├── handler_users_login.go         # POST /api/login
├── handler_users_update.go        # PUT /api/users
├── handler_users_refresh.go       # POST /api/refresh & /api/revoke
├── handler_webhooks.go            # POST /api/polka/webhooks
├── handler_readiness.go           # GET /api/healthz
├── handler_reset.go               # POST /admin/reset
├── metrics.go                     # GET /admin/metrics
├── utilities.go                   # JSON response helpers
├── internal/
│   ├── auth/
│   │   ├── auth.go                # JWT, Argon2id, token helpers
│   │   └── auth_test.go
│   └── database/                  # sqlc-generated DB layer
├── sql/
│   ├── queries/                   # SQL queries for sqlc
│   └── schema/                    # Goose migration files
├── sqlc.yaml
└── go.mod
```