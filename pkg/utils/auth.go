package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int, email, secret string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateTelegramAuth validates Telegram authentication data
func ValidateTelegramAuth(authData map[string]string, botToken string) error {
	// Check auth_date is within 24 hours
	authDate, ok := authData["auth_date"]
	if !ok {
		return fmt.Errorf("missing auth_date")
	}

	var authTimestamp int64
	fmt.Sscanf(authDate, "%d", &authTimestamp)
	if time.Now().Unix()-authTimestamp > 86400 {
		return fmt.Errorf("auth data is too old")
	}

	// Get hash and remove it from data
	hash, ok := authData["hash"]
	if !ok {
		return fmt.Errorf("missing hash")
	}

	// Create data check string
	var keys []string
	for key := range authData {
		if key != "hash" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var dataCheckArr []string
	for _, key := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", key, authData[key]))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// Create secret key
	secretKey := sha256.Sum256([]byte(botToken))

	// Calculate hash
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// Compare hashes
	if calculatedHash != hash {
		return fmt.Errorf("invalid hash")
	}

	return nil
}
