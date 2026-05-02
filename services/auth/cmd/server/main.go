package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	jwtadapter "github.com/harmeetsingh/grpc-envoy/services/auth/internal/adapter/jwt"
	httphandler "github.com/harmeetsingh/grpc-envoy/services/auth/internal/adapter/http"
	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/adapter/userstore"
	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	key := loadOrGenerateKey()

	store := userstore.NewHardcoded()
	signer := jwtadapter.NewRS256Signer(key, "https://auth.grpc-envoy.local")
	uc := usecase.New(store, signer)
	handler := httphandler.NewHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/login", handler.Login)

	log.Printf("auth-service listening on :%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadOrGenerateKey() *rsa.PrivateKey {
	keyPath := os.Getenv("RSA_KEY_PATH")
	if keyPath == "" {
		log.Println("RSA_KEY_PATH not set, generating ephemeral key")
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("failed to generate RSA key: %v", err)
		}
		return key
	}

	data, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("failed to read RSA key file %s: %v", keyPath, err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		log.Fatalf("failed to parse PEM block from %s", keyPath)
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse RSA private key: %v", err)
	}

	log.Printf("loaded RSA key from %s", keyPath)
	return key
}
