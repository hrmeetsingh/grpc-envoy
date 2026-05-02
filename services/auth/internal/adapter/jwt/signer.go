package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type RS256Signer struct {
	key    *rsa.PrivateKey
	issuer string
}

func NewRS256Signer(key *rsa.PrivateKey, issuer string) *RS256Signer {
	return &RS256Signer{key: key, issuer: issuer}
}

func (s *RS256Signer) Sign(userID string, tenant string) (string, error) {
	now := time.Now()
	claims := gojwt.MapClaims{
		"sub":    userID,
		"tenant": tenant,
		"iss":    s.issuer,
		"iat":    now.Unix(),
		"exp":    now.Add(15 * time.Minute).Unix(),
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodRS256, claims)
	return token.SignedString(s.key)
}

func (s *RS256Signer) Verify(tokenString string) (map[string]interface{}, error) {
	token, err := gojwt.Parse(tokenString, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return &s.key.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(gojwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
