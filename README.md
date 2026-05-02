# gRPC-Envoy

Go gRPC microservices with hexagonal architecture, fronted by an Envoy proxy that enables browser access via grpc-web. Includes a JWT auth service with subdomain-scoped tokens and proxy-level enforcement.

## Services

| Service | Port | Protocol | Methods | Description |
|---------|------|----------|---------|-------------|
| **UserService** | 50051 | gRPC | `CreateUser`, `GetUser` | In-memory user store |
| **GreeterService** | 50052 | gRPC | `SayHello`, `SayGoodbye` | Time-stamped greetings |
| **AuthService** | 8081 | HTTP | `POST /auth/login` | RS256 JWT with tenant claim |

## Architecture

Each gRPC service follows hexagonal (ports & adapters) architecture:

```
domain/     → entities, errors (zero dependencies)
usecase/    → application logic (depends on port interfaces)
port/       → interface definitions
adapter/
  grpcserver/ → inbound: proto → usecase
  memory/     → outbound: usecase → in-memory store
```

The auth service follows the same pattern with HTTP inbound adapter instead of gRPC.

See [docs/architecture.md](docs/architecture.md) for full diagram and details.

## How Envoy Works

Envoy sits between the browser and gRPC services, translating grpc-web (HTTP/1.1) to native gRPC (HTTP/2). It also enforces subdomain-token binding via a Lua filter.

See [docs/envoy.md](docs/envoy.md) for a detailed walkthrough of the config.

## Authentication & Subdomain Enforcement

The auth service issues JWTs with a `tenant` claim derived from the request's `Host` header subdomain:

```
POST /auth/login → Host: acme.example.com
Response: { "token": "<JWT with tenant=acme>" }
```

On all subsequent gRPC requests, Envoy's Lua filter (`envoy/filters/tenant_check.lua`):
1. Extracts the `tenant` claim from the Bearer token
2. Extracts the subdomain from the request's `Host` header
3. Returns **403** if they don't match

This prevents users from changing subdomains in the URL to access other tenants' data.

### Test users (hardcoded)

| Username | Password |
|----------|----------|
| alice | password123 |
| bob | password456 |

### Testing with curl

Start the stack with `make up`, then try these commands:

**1. Login and get a token (scoped to `acme` subdomain):**

```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H "Host: acme.example.com" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'
```

Response:
```json
{"token":"eyJhbGciOiJSUzI1NiIs..."}
```

Save the token for subsequent requests:

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Host: acme.example.com" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
```

**2. Request with correct subdomain — passes through:**

```bash
curl -s -w "\nHTTP %{http_code}\n" \
  http://localhost:8080/user.v1.UserService/CreateUser \
  -H "Host: acme.example.com" \
  -H "Authorization: Bearer $TOKEN" \
  --max-time 2
```

This passes the Lua filter (the request reaches the upstream gRPC service).

**3. Request with wrong subdomain — blocked with 403:**

```bash
curl -s -w "\nHTTP %{http_code}\n" \
  http://localhost:8080/user.v1.UserService/CreateUser \
  -H "Host: evil.example.com" \
  -H "Authorization: Bearer $TOKEN" \
  --max-time 2
```

Response:
```json
{"error":"token not valid for this subdomain"}
HTTP 403
```

The token was issued for `acme` but the request targets `evil` — Envoy rejects it.

**4. Request with no auth header — rejected with 401:**

```bash
curl -s -w "\nHTTP %{http_code}\n" \
  http://localhost:8080/user.v1.UserService/CreateUser \
  -H "Host: acme.example.com" \
  --max-time 2
```

Response:
```json
{"error":"missing authorization header"}
HTTP 401
```

**5. Login with invalid credentials — rejected:**

```bash
curl -s -w "\nHTTP %{http_code}\n" -X POST \
  http://localhost:8080/auth/login \
  -H "Host: acme.example.com" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"wrong"}'
```

Response:
```json
{"error":"invalid credentials"}
HTTP 401
```

### RSA Key Configuration

By default, the auth service generates an ephemeral RSA key on startup. For persistence across restarts, set:

```bash
RSA_KEY_PATH=/path/to/private-key.pem
```

## Prerequisites

- Go 1.25+
- Node.js 22+
- Docker & Docker Compose
- protoc with `protoc-gen-go` and `protoc-gen-go-grpc` (for regenerating stubs)

## Quick Start

### Run with Docker Compose

```bash
make up
```

This starts all five containers:
- **user-service** on `:50051`
- **greeter-service** on `:50052`
- **auth-service** on `:8081`
- **envoy** on `:8080` (proxy) and `:9901` (admin)
- **frontend** on `:3000`

Open http://localhost:3000 in your browser.

### Stop

```bash
make down
```

## Development

### Run tests

```bash
make test
```

### Build binaries locally

```bash
make build
```

Outputs to `bin/user-service`, `bin/greeter-service`, and `bin/auth-service`.

### Regenerate protobuf stubs

```bash
make proto
```

### Run services locally (without Docker)

Terminal 1:
```bash
cd services/user && go run ./cmd/server
```

Terminal 2:
```bash
cd services/greeter && go run ./cmd/server
```

Terminal 3:
```bash
cd services/auth && go run ./cmd/server
```

Terminal 4 (requires envoy installed):
```bash
envoy -c envoy/envoy.yaml
```

Terminal 5:
```bash
cd frontend && npm run dev
```

## Project Structure

```
grpc-envoy/
├── proto/                       # Protobuf definitions
│   ├── user/v1/user.proto
│   ├── greeter/v1/greeter.proto
│   ├── buf.yaml
│   └── buf.gen.yaml
├── services/
│   ├── user/                    # UserService (hexagonal, gRPC)
│   │   ├── cmd/server/main.go
│   │   ├── internal/
│   │   └── Dockerfile
│   ├── greeter/                 # GreeterService (mirrors user/)
│   └── auth/                    # AuthService (hexagonal, HTTP)
│       ├── cmd/server/main.go
│       ├── internal/
│       │   ├── domain/          # Credentials, TokenResult, errors
│       │   ├── usecase/         # Login logic + subdomain extraction
│       │   ├── port/            # UserStore, TokenSigner interfaces
│       │   └── adapter/
│       │       ├── http/        # HTTP handler (POST /auth/login)
│       │       ├── jwt/         # RS256 signer/verifier
│       │       └── userstore/   # Hardcoded test users
│       └── Dockerfile
├── envoy/
│   ├── envoy.yaml               # Envoy config (grpc_web + Lua + CORS)
│   └── filters/
│       └── tenant_check.lua     # Subdomain enforcement filter
├── frontend/                    # Next.js (App Router, TS, Tailwind)
├── docker-compose.yml
├── Makefile
└── docs/
    ├── architecture.md          # System diagram
    └── envoy.md                 # Envoy proxy documentation
```
