package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT JWT 토큰 생성
func GenerateJWT(userID, email, role, secret string, expiresHours int) (string, error) {
	claims := &JWTClaim{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiresHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT JWT 토큰 검증
func ValidateJWT(tokenString string, secret string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		// 서명 방식 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GetJWTClaims HTTP 요청에서 JWT 토큰을 추출하고 검증
func GetJWTClaims(r *http.Request, secret string) (*JWTClaim, error) {
	// Authorization 헤더에서 토큰 추출
	tokenString := extractToken(r)
	if tokenString == "" {
		return nil, fmt.Errorf("no token found")
	}

	return ValidateJWT(tokenString, secret)
}

// extractToken Authorization 헤더에서 토큰 추출
func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 && strArr[0] == "Bearer" {
		return strArr[1]
	}
	return ""
}
