# Envoy Proxy — How It Works

## The Problem

Browsers cannot speak native gRPC. gRPC uses HTTP/2 with binary protobuf framing, trailers, and bidirectional streaming — features that browser `fetch`/`XMLHttpRequest` APIs don't support directly.

## The Solution: grpc-web

The [grpc-web](https://github.com/grpc/grpc-web) protocol is a browser-compatible variant of gRPC. It wraps protobuf payloads in a slightly different framing that works over HTTP/1.1. A proxy sits between the browser and the gRPC server, translating between grpc-web (HTTP/1.1) and native gRPC (HTTP/2).

Envoy is the officially recommended proxy for this translation.

## How Envoy Is Configured

The configuration lives in `envoy/envoy.yaml`. Here's a breakdown of each major section:

### Admin Interface

```yaml
admin:
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 9901
```

Exposes Envoy's built-in admin dashboard at `http://localhost:9901`. Useful for inspecting clusters, routes, stats, and config dumps during development.

### Listener

```yaml
listeners:
  - name: listener_0
    address:
      socket_address:
        address: 0.0.0.0
        port_value: 8080
```

Envoy listens on port 8080 for incoming HTTP requests from the browser.

### HTTP Connection Manager

The listener uses the `http_connection_manager` network filter, which handles HTTP-level processing:

- **codec_type: AUTO** — auto-detect HTTP/1.1 or HTTP/2
- **route_config** — defines how incoming requests are routed to upstream gRPC clusters

### Routing

```yaml
routes:
  - match:
      prefix: "/user.v1.UserService"
    route:
      cluster: user_service
  - match:
      prefix: "/greeter.v1.GreeterService"
    route:
      cluster: greeter_service
```

gRPC requests use the path format `/<package>.<Service>/<Method>`. Envoy matches on the service prefix and forwards to the correct upstream cluster.

### HTTP Filters (the key part)

Three filters execute in order on every request:

1. **`grpc_web`** — Translates between grpc-web wire format (what the browser sends) and standard gRPC (what the Go servers expect). This is the core filter that makes browser-to-gRPC possible.

2. **`cors`** — Handles Cross-Origin Resource Sharing. The browser's frontend (port 3000) makes requests to Envoy (port 8080), which is a different origin. Without CORS headers, the browser would block the request.

3. **`router`** — The standard Envoy HTTP router that forwards requests to upstream clusters based on the route config.

### CORS Configuration

```yaml
cors:
  allow_origin_string_match:
    - prefix: "*"
  allow_methods: GET, PUT, DELETE, POST, OPTIONS
  allow_headers: keep-alive,user-agent,cache-control,content-type,...,x-grpc-web,grpc-timeout
  expose_headers: grpc-status,grpc-message
```

- `allow_origin_string_match: *` — accepts requests from any origin (fine for development)
- `allow_headers` — must include `x-grpc-web`, `content-type`, and `grpc-timeout` for grpc-web to work
- `expose_headers` — exposes `grpc-status` and `grpc-message` so the browser client can read gRPC error details

### Upstream Clusters

```yaml
clusters:
  - name: user_service
    type: LOGICAL_DNS
    typed_extension_protocol_options:
      envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
        explicit_http_config:
          http2_protocol_options: {}
    load_assignment:
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: user-service
                    port_value: 50051
```

Each cluster defines one upstream gRPC service:

- **type: LOGICAL_DNS** — resolves hostnames via DNS (docker-compose service names work here)
- **http2_protocol_options** — forces HTTP/2 to the upstream, which gRPC requires
- **address** — the docker-compose service name (`user-service` or `greeter-service`)

## Request Flow Diagram

```
Browser (port 3000)
    │
    │  HTTP/1.1 POST /user.v1.UserService/CreateUser
    │  Content-Type: application/grpc-web+proto
    │  X-Grpc-Web: 1
    │
    ▼
Envoy (port 8080)
    │
    │  [grpc_web filter] strips grpc-web framing, converts to gRPC
    │  [cors filter] adds Access-Control-Allow-* headers
    │  [router] matches prefix → user_service cluster
    │
    │  HTTP/2 POST /user.v1.UserService/CreateUser
    │  Content-Type: application/grpc
    │
    ▼
UserService (port 50051)
    │
    │  Processes request, returns protobuf response
    │
    ▼
Envoy
    │
    │  [grpc_web filter] re-frames response for browser
    │
    ▼
Browser
    │
    │  grpc-web client deserializes protobuf → JavaScript object
```

## Key Takeaways

1. **Envoy is a transparent bridge** — your Go services don't know or care that the client is a browser. They receive standard gRPC requests.
2. **The grpc_web filter does the heavy lifting** — it translates between HTTP/1.1 grpc-web format and HTTP/2 native gRPC.
3. **CORS is required** — browser security model requires explicit permission for cross-origin requests.
4. **HTTP/2 to upstream is mandatory** — gRPC requires HTTP/2, so `http2_protocol_options` must be set on every gRPC cluster.
5. **Routing uses gRPC path convention** — `/<package>.<Service>/<Method>`, and Envoy matches on the service prefix.
