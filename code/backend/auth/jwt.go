package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"digital-innovation/stratego/models"
	"digital-innovation/stratego/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var jwtSecret = []byte(utils.GetEnv("JWT_SECRET", "dev_secret_only_for_local_development_123"))

type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JWTPayload struct {
	Sub  int    `json:"sub"`
	Name string `json:"name"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
}

func GenerateToken(userID int, username string) (string, error) {
	header := JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	now := time.Now()
	payload := JWTPayload{
		Sub:  userID,
		Name: username,
		Iat:  now.Unix(),
		Exp:  now.Add(30 * 24 * time.Hour).Unix(), // 30 days
	}

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	unsignedToken := headerBase64 + "." + payloadBase64
	signature := computeHmac(unsignedToken, jwtSecret)

	return unsignedToken + "." + signature, nil
}

func VerifyToken(tokenString string) (*models.User, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerBase64, payloadBase64, signature := parts[0], parts[1], parts[2]
	unsignedToken := headerBase64 + "." + payloadBase64

	expectedSignature := computeHmac(unsignedToken, jwtSecret)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, errors.New("invalid signature")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var payload JWTPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if time.Now().Unix() > payload.Exp {
		return nil, errors.New("token expired")
	}

	return &models.User{
		ID:       payload.Sub,
		Username: payload.Name,
	}, nil
}

func computeHmac(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
