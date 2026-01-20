package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"liangxiong/demo/internal/config"
)

// JWTManager encapsulates signing and verification logic.
type JWTManager struct {
	issuer   string
	audience string
	ttl      time.Duration
	method   jwt.SigningMethod
	secret   []byte
	privKey  *rsa.PrivateKey
	pubKey   *rsa.PublicKey
}

// NewJWTManager builds a manager using the provided configuration.
func NewJWTManager(cfg config.AuthConfig) (*JWTManager, error) {
	manager := &JWTManager{
		issuer:   cfg.Issuer,
		audience: cfg.Audience,
		ttl:      cfg.AccessTokenTTL,
	}

	switch strings.ToUpper(cfg.Algorithm) {
	case "HS256":
		manager.method = jwt.SigningMethodHS256
		manager.secret = []byte(cfg.JWTSecret)
	case "RS256":
		manager.method = jwt.SigningMethodRS256
		priv, err := readPrivateKey(cfg.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
		pub, err := readPublicKey(cfg.PublicKeyPath)
		if err != nil {
			return nil, err
		}
		manager.privKey = priv
		manager.pubKey = pub
	default:
		return nil, fmt.Errorf("unsupported algorithm %s", cfg.Algorithm)
	}

	return manager, nil
}

// Generate issues a signed JWT for the subject.
func (m *JWTManager) Generate(userID string) (string, time.Time, error) {
	if m == nil {
		return "", time.Time{}, errors.New("jwt manager is nil")
	}

	expiresAt := time.Now().Add(m.ttl)
	claims := jwt.RegisteredClaims{
		Issuer:    m.issuer,
		Audience:  jwt.ClaimStrings{m.audience},
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(m.method, claims)
	signed, err := token.SignedString(m.signingKey())
	return signed, expiresAt, err
}

// Validate parses the token string and returns claims if valid.
func (m *JWTManager) Validate(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != m.method.Alg() {
			return nil, fmt.Errorf("unexpected jwt signing method %s", token.Header["alg"])
		}
		return m.verificationKey(), nil
	}, jwt.WithAudience(m.audience), jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (m *JWTManager) signingKey() interface{} {
	if m.method == jwt.SigningMethodHS256 {
		return m.secret
	}
	return m.privKey
}

func (m *JWTManager) verificationKey() interface{} {
	if m.method == jwt.SigningMethodHS256 {
		return m.secret
	}
	return m.pubKey
}

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid private key pem")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsed, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		if pk, ok := parsed.(*rsa.PrivateKey); ok {
			return pk, nil
		}
		return nil, errors.New("pkcs8 key is not rsa")
	}
	return key, nil
}

func readPublicKey(path string) (*rsa.PublicKey, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid public key pem")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not rsa")
	}
	return rsaPub, nil
}
