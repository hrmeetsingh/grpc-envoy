# Architecture

## System Diagram

```mermaid
flowchart LR
  subgraph browser["Browser"]
    nextApp["Next.js App<br/>(App Router + TS)"]
    grpcWebClient["grpc-web client<br/>(generated TS stubs)"]
  end

  subgraph edge["Edge"]
    envoy["Envoy Proxy<br/>:8080<br/>grpc_web filter + CORS"]
  end

  subgraph services["gRPC Services (Go, hexagonal)"]
    subgraph userSvc["UserService :50051"]
      userAdapter["gRPC adapter<br/>(inbound)"]
      userCore["domain + usecase<br/>CreateUser / GetUser"]
      userRepo["in-memory repo<br/>(outbound adapter)"]
    end
    subgraph greeterSvc["GreeterService :50052"]
      greeterAdapter["gRPC adapter<br/>(inbound)"]
      greeterCore["domain + usecase<br/>SayHello / SayGoodbye"]
      greeterClock["clock port<br/>(outbound adapter)"]
    end
  end

  protoFiles["proto/<br/>user.proto<br/>greeter.proto"]

  nextApp --> grpcWebClient
  grpcWebClient -- "HTTP/1.1 grpc-web" --> envoy
  envoy -- "HTTP/2 gRPC" --> userAdapter
  envoy -- "HTTP/2 gRPC" --> greeterAdapter
  userAdapter --> userCore --> userRepo
  greeterAdapter --> greeterCore --> greeterClock
  protoFiles -. "protoc generate" .-> userAdapter
  protoFiles -. "protoc generate" .-> greeterAdapter
  protoFiles -. "protoc generate" .-> grpcWebClient
```

## Hexagonal Architecture

Each Go service follows the ports-and-adapters (hexagonal) pattern:

```
services/<name>/
├── cmd/server/main.go          # Composition root — wires adapters to ports
└── internal/
    ├── domain/                 # Core entities and domain errors
    ├── usecase/                # Application logic (depends only on ports)
    ├── port/                   # Interfaces (inbound + outbound contracts)
    └── adapter/
        ├── grpcserver/         # Inbound adapter: gRPC → usecase
        └── memory/             # Outbound adapter: usecase → in-memory store
```

**Key rules:**
- Domain and usecase packages have zero infrastructure imports.
- Adapters depend inward on ports; ports never depend on adapters.
- `main.go` is the only place that knows about concrete adapter types.

## Services

| Service | Port | Methods |
|---------|------|---------|
| UserService | 50051 | `CreateUser(name) → (id, name)`, `GetUser(id) → (id, name)` |
| GreeterService | 50052 | `SayHello(name) → message`, `SayGoodbye(name) → message` |

## Data Flow

1. Browser makes HTTP/1.1 request using grpc-web binary format
2. Envoy's `grpc_web` filter translates to HTTP/2 gRPC
3. Envoy routes by service prefix (`/user.v1.UserService`, `/greeter.v1.GreeterService`)
4. Go gRPC adapter deserializes proto, calls usecase
5. Usecase executes domain logic via port interfaces
6. Response flows back through the same chain
