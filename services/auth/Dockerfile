FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /auth-service ./cmd/server

FROM alpine:3.20
COPY --from=build /auth-service /auth-service
EXPOSE 8081
ENTRYPOINT ["/auth-service"]
