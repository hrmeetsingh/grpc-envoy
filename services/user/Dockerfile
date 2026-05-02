FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /user-service ./cmd/server

FROM alpine:3.20
COPY --from=build /user-service /user-service
EXPOSE 50051
ENTRYPOINT ["/user-service"]
