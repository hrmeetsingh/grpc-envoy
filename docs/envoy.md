# Envoy Proxy — How It Works

## The Problem

Browsers cannot speak native gRPC. gRPC uses HTTP/2 with protobuf framing, trailers, and streaming semantics that typical browser APIs do not expose cleanly. This project also routes **HTTP login** and **gRPC** through one listener and enforces **JWT tenant vs subdomain** before upstream gRPC.

## grpc-web and Envoy

The [grpc-web](https://github.com/grpc/grpc-web) protocol carries protobuf payloads over HTTP in a form browsers can use. Envoy sits between the browser and Go gRPC servers and translates grpc-web (HTTP/1.1 as seen on the client side) to native gRPC (HTTP/2) upstream.

## Configuration Location

- Main config: `envoy/envoy.yaml`
- Lua tenant script (mounted in Docker): `envoy/filters/tenant_check.lua`

## Admin Interface

```yaml
admin:
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 9901
```

Dashboard: `http://localhost:9901` (clusters, routes, stats).

## Listener

Envoy listens on **0.0.0.0:8080** for all traffic from the browser and for curl tests against `localhost:8080`.

## HTTP Connection Manager

- **codec_type: AUTO** — HTTP/1.1 or HTTP/2 on the downstream connection.
- **route_config** — path-prefix routes to `auth_service`, `user_service`, or `greeter_service`.

## Routing

Routes are **order-sensitive** (first match wins):

```yaml
routes:
  - match:
      prefix: "/auth/"
    route:
      cluster: auth_service
  - match:
      prefix: "/user.v1.UserService"
    route:
      cluster: user_service
  - match:
      prefix: "/greeter.v1.GreeterService"
    route:
      cluster: greeter_service
```

- **Auth** — JSON login and any path under `/auth/` goes to Go **AuthService** on port **8081** (plain HTTP).
- **gRPC** — paths follow `/<package>.<Service>/<Method>`; clusters use **HTTP/2** upstream.

## HTTP Filters

Execution order (request path):

1. **`envoy.filters.http.grpc_web`** — grpc-web to native gRPC toward upstream (login traffic is plain HTTP; filter still runs but routing targets HTTP/1 auth cluster).
2. **`envoy.filters.http.cors`** — Adds CORS headers per virtual host (`cors` block on the route table).
3. **`envoy.filters.http.lua`** — Loads `default_source_code.filename: /etc/envoy/filters/tenant_check.lua`. For paths **not** starting with `/auth/`, requires `Authorization: Bearer`, decodes JWT payload enough to read `tenant`, compares with subdomain from `:authority`; responds **401** or **403** on failure; **skips** checks for `/auth/*` so login works without a prior token.
4. **`envoy.filters.http.router`** — Selects cluster from route config.

In Docker Compose, `./envoy/filters` is mounted at `/etc/envoy/filters` so the Lua file can be edited without baking it into an image layer as inline YAML.

### CORS

```yaml
cors:
  allow_origin_string_match:
    - prefix: "*"
  allow_methods: GET, PUT, DELETE, POST, OPTIONS
  allow_headers: ... ,x-grpc-web,grpc-timeout,authorization
```

- **`authorization`** — required so browsers can send `Bearer` tokens after login.
- **grpc-web** — still needs `x-grpc-web`, `content-type`, `grpc-timeout` as before.

## Upstream Clusters

**User and Greeter** (gRPC):

```yaml
clusters:
  - name: user_service
    typed_extension_protocol_options:
      envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
        explicit_http_config:
          http2_protocol_options: {}
```

- **`http2_protocol_options`** — Required for native gRPC to Go servers.

**Auth** (HTTP only):

```yaml
clusters:
  - name: auth_service
    load_assignment:
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: auth-service
                    port_value: 8081
```

No `http2_protocol_options`: auth speaks normal HTTP/1.1.

- **LOGICAL_DNS** — docker-compose DNS names (`user-service`, `greeter-service`, `auth-service`).

## Request Flow Diagrams

### Login (`POST /auth/login`)

```
Browser or curl → Envoy :8080
  Host: tenant.example.com
  Path /auth/login
    → grpc_web (no-op for JSON body upstream)
    → cors
    → lua (skip — path /auth/)
    → router → auth_service :8081
    → JSON { "token": "..." }
```

### gRPC call (after login)

```
Browser → Envoy :8080
  Authorization: Bearer <JWT>
  Path /user.v1.UserService/CreateUser
    → grpc_web (grpc-web → gRPC)
    → cors
    → lua (tenant vs Host subdomain; fail fast 401/403)
    → router → user_service :50051
    → HTTP/2 gRPC to Go UserService
```

Response path mirrors filters in reverse order; grpc_web reframes grpc-web for the browser where applicable.

## Key Takeaways

1. **One listener** multiplexes REST login and gRPC by **path prefix**.
2. **`grpc_web`** bridges browser grpc-web and upstream gRPC (**HTTP/2** to user/greeter).
3. **Lua** enforces **tenant ⊂ JWT** aligns with **`Host` subdomain** for all non-login HTTP requests that reach the filter chain (including grpc-web-encoded calls).
4. **CORS** must allow **`authorization`** header for Bearer tokens.
5. **auth_service cluster** stays **HTTP/1**; **gRPC clusters** require **HTTP/2** upstream options.
