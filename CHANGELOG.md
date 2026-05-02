# Changelog

## v1 - Clarify - 2026-05-01T12:27

### What changed
- Captured project requirements and locked defaults (user skipped clarifying questions)

### Files touched
- `CHANGELOG.md`

### Test status
- no tests yet

### Requirements
- Two gRPC services with hexagonal architecture
- Each service exposes 2 methods (minimal functionality, future-extensible)
- Envoy proxy bridges browser HTTP clients to gRPC backends
- Next.js frontend consumes endpoints through Envoy
- Minimum docs sufficient to understand Envoy proxy role
- Architecture diagram saved to a markdown file

### Defaults (no answers received, sensible picks)
- **Language**: Go (workspace path `go/grpc-envoy`)
- **Service A**: `UserService` — `CreateUser`, `GetUser`
- **Service B**: `GreeterService` — `SayHello`, `SayGoodbye`
- **Storage**: in-memory adapter behind repository port
- **Envoy mode**: `envoy.filters.http.grpc_web` (browser uses grpc-web client)
- **Orchestration**: docker-compose + Makefile shortcuts
- **Frontend**: Next.js App Router, TypeScript, Tailwind, `grpc-web` JS client
- **Test scope**: domain unit tests + gRPC adapter integration tests

## v2 - Diagram - 2026-05-01T12:51

### What changed
- Architecture diagram presented and confirmed
- Directory layout agreed upon

### Files touched
- `CHANGELOG.md`
- `docs/architecture.md` (planned)

### Test status
- no tests yet

## v3 - Tests First - 2026-05-01T12:52

### What changed
- Created domain entities, port interfaces, and stub implementations for both services
- Wrote 8 usecase tests (4 user, 4 greeter) + 3 adapter tests (2 user memory repo, 1 greeter clock)
- All tests fail with `panic: not implemented` as expected

### Files touched
- `services/user/internal/domain/user.go`
- `services/user/internal/port/repository.go`
- `services/user/internal/usecase/user_usecase.go`
- `services/user/internal/usecase/user_usecase_test.go`
- `services/user/internal/adapter/memory/repository.go`
- `services/user/internal/adapter/memory/repository_test.go`
- `services/greeter/internal/domain/greeting.go`
- `services/greeter/internal/port/clock.go`
- `services/greeter/internal/usecase/greeter_usecase.go`
- `services/greeter/internal/usecase/greeter_usecase_test.go`
- `services/greeter/internal/adapter/memory/clock.go`
- `services/greeter/internal/adapter/memory/clock_test.go`

### Test status
- 0 passing / 10 total (all panic: not implemented)

## v4 - Minimal Implementation - 2026-05-01T12:55

### What changed
- Implemented usecase and adapter code for both services — all 10 tests green
- Generated protobuf Go stubs via protoc for user.proto and greeter.proto
- Built gRPC inbound adapters (grpcserver/server.go) mapping proto ↔ domain
- Created main.go binaries for both services (ports 50051, 50052)
- Configured Envoy proxy (grpc_web filter, CORS, two upstream clusters)
- docker-compose.yml orchestrating user-service, greeter-service, envoy, frontend
- Makefile with proto, test, build, up, down targets
- Scaffolded Next.js frontend (App Router, TS, Tailwind) with grpc-web client stubs
- Frontend builds successfully (Next.js 14 optimized production build)

### Files touched
- `services/user/internal/usecase/user_usecase.go` — implemented CreateUser, GetUser
- `services/user/internal/adapter/memory/repository.go` — concurrent-safe in-memory store
- `services/greeter/internal/usecase/greeter_usecase.go` — implemented SayHello, SayGoodbye
- `services/greeter/internal/adapter/memory/clock.go` — real clock adapter
- `proto/user/v1/user.proto`, `proto/greeter/v1/greeter.proto`
- `proto/buf.yaml`, `proto/buf.gen.yaml`
- `services/*/internal/adapter/grpcserver/pb/**` — generated pb files
- `services/*/internal/adapter/grpcserver/server.go` — gRPC adapters
- `services/*/cmd/server/main.go` — service entrypoints
- `services/*/Dockerfile`
- `envoy/envoy.yaml`
- `docker-compose.yml`
- `Makefile`
- `frontend/` — full Next.js app with grpc-web client stubs

### Test status
- 10 passing / 10 total

### Summary
Complete vertical slice: proto → generated stubs → hexagonal Go services (domain, usecase, ports, adapters) → Envoy grpc-web proxy → Next.js browser client. All domain and adapter tests pass. Frontend compiles and produces optimized build. Stack runs via docker-compose.

## v5 - Refactor Plan - 2026-05-01T13:12

### What changed
- Proposed 5-point refactor plan, confirmed by user

### Files touched
- `CHANGELOG.md`

### Test status
- 10 passing / 10 total

### Plan
1. Add `.gitignore` at repo root
2. Ensure consistent `toStatus` error mapping style in both gRPC adapters
3. Pass `context.Context` through consistently in greeter usecase
4. Remove unused connect-es plugin from buf.gen.yaml
5. Add envoy admin interface on port 9901

## v6 - Refactor - 2026-05-01T13:14

### What changed
- Added `.gitignore` at repo root (bin/, node_modules/, .next/, build artifacts)
- Greeter usecase now accepts and checks `context.Context` (cancellation guard)
- Removed unused `buf.build/connectrpc/es` plugin from `buf.gen.yaml`
- Added Envoy admin interface on port 9901 in `envoy.yaml` and exposed in docker-compose
- User and greeter gRPC adapters already had consistent `toStatus` style — no change needed

### Files touched
- `.gitignore` (new)
- `services/greeter/internal/usecase/greeter_usecase.go`
- `proto/buf.gen.yaml`
- `envoy/envoy.yaml`
- `docker-compose.yml`

### Test status
- 10 passing / 10 total

### Summary
Minor cleanups: proper gitignore, context propagation, remove dead config, add envoy admin for debugging. No structural changes. All tests remain green.

## v7 - README - 2026-05-01T13:16

### What changed
- Created `README.md` with project overview, quick start, dev instructions, project structure
- Created `docs/architecture.md` with mermaid system diagram, hexagonal architecture explanation, service table, data flow
- Created `docs/envoy.md` with detailed Envoy proxy documentation: problem statement, config breakdown (admin, listener, routing, filters, CORS, clusters), request flow diagram, key takeaways

### Files touched
- `README.md` (new)
- `docs/architecture.md` (new)
- `docs/envoy.md` (new)

### Test status
- 10 passing / 10 total

## v8 - Clarify (Auth Service) - 2026-05-01T23:11

### What changed
- Captured requirements for new auth service feature

### Files touched
- `CHANGELOG.md`

### Test status
- 10 passing / 10 total (existing, no new tests yet)

### Requirements
- **Language**: Go (HTTP/REST on port 8081)
- **JWT signing**: RS256 (asymmetric key pair)
- **User store**: Hardcoded test users
- **Subdomain source**: Extract from Host header of login request
- **Endpoints**: POST /login → JWT
- **Envoy enforcement**: Lua filter compares jwt.tenant claim to request subdomain

## v9 - Diagram (Auth Service) - 2026-05-01T23:13

### What changed
- Architecture diagram presented and confirmed
- Auth service routes through Envoy on /auth/* prefix
- Lua filter inserted into HTTP filter chain before router
- JWT contains tenant claim derived from Host header subdomain

### Files touched
- `CHANGELOG.md`

### Test status
- 10 passing / 10 total (existing)

## v10 - Tests First (Auth Service) - 2026-05-01T23:15

### What changed
- Created auth service skeleton: domain, ports, stub adapters
- Wrote 15 failing tests across 3 packages:
  - usecase (6): login valid/invalid, tenant extraction, empty host, host with port
  - jwt adapter (6): sign, verify, claims content, wrong key, invalid token
  - userstore adapter (3): valid user, invalid password, unknown user
- All fail with `panic: not implemented`

### Files touched
- `services/auth/go.mod`
- `services/auth/internal/domain/auth.go`
- `services/auth/internal/port/userstore.go`
- `services/auth/internal/port/signer.go`
- `services/auth/internal/usecase/auth_usecase.go`
- `services/auth/internal/usecase/auth_usecase_test.go`
- `services/auth/internal/adapter/jwt/signer.go`
- `services/auth/internal/adapter/jwt/signer_test.go`
- `services/auth/internal/adapter/userstore/hardcoded.go`
- `services/auth/internal/adapter/userstore/hardcoded_test.go`

### Test status
- 0 passing / 15 total (all panic: not implemented)

## v11 - Minimal Implementation (Auth Service) - 2026-05-01T23:20

### What changed
- Implemented RS256 JWT signer (sign + verify with golang-jwt/jwt/v5)
- Implemented hardcoded user store (alice/password123, bob/password456)
- Implemented auth usecase: validates credentials, extracts subdomain from Host, issues JWT with tenant claim
- Created HTTP handler: POST /auth/login → JSON {token}
- Created main.go for auth-service on port 8081
- Created Dockerfile for auth-service
- Updated docker-compose.yml: added auth-service container
- Updated envoy.yaml:
  - Added /auth/* route to auth_service cluster
  - Added Lua HTTP filter: extracts JWT tenant claim, compares to Host subdomain
  - Skips auth check for /auth/ paths (login doesn't require token)
  - Returns 401 for missing/invalid Authorization header
  - Returns 403 when tenant claim doesn't match request subdomain
  - Added `authorization` to CORS allowed headers
  - Added auth_service cluster (HTTP/1.1, port 8081)
- Updated Makefile: added auth service to test and build targets
- Verified end-to-end: login returns JWT, wrong subdomain → 403, no auth → 401

### Files touched
- `services/auth/internal/adapter/jwt/signer.go`
- `services/auth/internal/adapter/userstore/hardcoded.go`
- `services/auth/internal/usecase/auth_usecase.go`
- `services/auth/internal/adapter/http/handler.go`
- `services/auth/cmd/server/main.go`
- `services/auth/Dockerfile`
- `services/auth/go.mod`, `services/auth/go.sum`
- `docker-compose.yml`
- `envoy/envoy.yaml`
- `Makefile`

### Test status
- 25 passing / 25 total (10 existing + 15 auth)

### Summary
Complete auth service: hardcoded users, RS256 JWT with tenant claim derived from Host header subdomain, HTTP login endpoint behind Envoy. Envoy Lua filter enforces subdomain-token binding on all gRPC requests — blocks cross-subdomain access with 403. Auth endpoints (/auth/*) bypass the check.

## v12 - Refactor Plan (Auth Service) - 2026-05-01T23:23

### What changed
- Proposed 5-point refactor plan, confirmed by user

### Files touched
- `CHANGELOG.md`

### Test status
- 25 passing / 25 total

### Plan
1. Extract Lua filter to separate file (envoy/filters/tenant_check.lua)
2. Add HTTP handler unit tests
3. Make RSA key configurable (load from file path)
4. Add context cancellation check in auth usecase
5. Validate login request fields (empty username/password → 400)

## v13 - Refactor (Auth Service) - 2026-05-01T23:25

### What changed
1. Extracted Lua filter to `envoy/filters/tenant_check.lua`, referenced via `default_source_code.filename`
2. Added 5 HTTP handler tests (method not allowed, invalid JSON, bad creds, success, empty fields)
3. RSA key now loaded from `RSA_KEY_PATH` env var; falls back to ephemeral key if unset
4. Auth usecase checks `ctx.Err()` before proceeding
5. Handler validates empty username/password → returns 400 before hitting usecase
6. Mounted `envoy/filters/` volume in docker-compose for Envoy

### Files touched
- `envoy/filters/tenant_check.lua` (new)
- `envoy/envoy.yaml`
- `docker-compose.yml`
- `services/auth/cmd/server/main.go`
- `services/auth/internal/usecase/auth_usecase.go`
- `services/auth/internal/adapter/http/handler.go`
- `services/auth/internal/adapter/http/handler_test.go` (new)

### Test status
- 30 passing / 30 total (10 user + 10 greeter + 20 auth)

### Summary
Refactored for maintainability: Lua filter lives in its own file for easier editing and version control. RSA key is configurable for persistence across restarts. Input validation moved to handler layer. Added comprehensive handler tests.

## v14 - README - 2026-05-01T23:28

### What changed
- Updated README.md with auth service documentation
- Added service table entry for AuthService
- Documented authentication flow and subdomain enforcement
- Listed test users and RSA key configuration
- Updated quick start (5 containers), project structure, local dev instructions

### Files touched
- `README.md`

### Test status
- 30 passing / 30 total
