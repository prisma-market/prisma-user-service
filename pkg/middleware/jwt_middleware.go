package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/kihyun1998/prisma-market/prisma-user-service/pkg/utils"
)

type JWTMiddleware struct {
	jwtSecret string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{
		jwtSecret: secret,
	}
}

// ValidateJWT HTTP 요청에 대한 JWT 검증 미들웨어
func (m *JWTMiddleware) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// JWT 토큰 검증
		claims, err := utils.GetJWTClaims(r, m.jwtSecret)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "unauthorized: " + err.Error(),
			})
			return
		}

		// 컨텍스트에 claims 정보 저장
		ctx := utils.SetUserContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole 특정 역할이 필요한 엔드포인트를 위한 미들웨어
func (m *JWTMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := utils.GetJWTClaims(r, m.jwtSecret)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(ErrorResponse{
					Error: "unauthorized: " + err.Error(),
				})
				return
			}

			// 역할 확인
			if claims.Role != role {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(ErrorResponse{
					Error: "forbidden: insufficient permissions",
				})
				return
			}

			ctx := utils.SetUserContext(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
