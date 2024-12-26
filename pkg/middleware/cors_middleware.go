package middleware

import (
	"net/http"
	"strings"
)

type CORS struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
	maxAge         int // preflight 요청 캐시 시간(초)
}

func NewCORS() *CORS {
	return &CORS{
		allowedOrigins: []string{"*"}, // 실제 운영 환경에서는 구체적인 도메인 지정 필요
		allowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		allowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Request-ID",
		},
		maxAge: 86400, // 24시간
	}
}

// WithOrigins CORS 허용 도메인 설정
func (c *CORS) WithOrigins(origins ...string) *CORS {
	c.allowedOrigins = origins
	return c
}

// WithMethods CORS 허용 메서드 설정
func (c *CORS) WithMethods(methods ...string) *CORS {
	c.allowedMethods = methods
	return c
}

// WithHeaders CORS 허용 헤더 설정
func (c *CORS) WithHeaders(headers ...string) *CORS {
	c.allowedHeaders = headers
	return c
}

// WithMaxAge preflight 요청 캐시 시간 설정
func (c *CORS) WithMaxAge(seconds int) *CORS {
	c.maxAge = seconds
	return c
}

// Handler CORS 미들웨어 핸들러
func (c *CORS) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Origin 검사
		if origin != "" {
			if !c.isAllowedOrigin(origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Preflight 요청 처리
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.allowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.allowedHeaders, ","))
			w.Header().Set("Access-Control-Max-Age", string(c.maxAge))
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Vary 헤더 설정
			// 캐시된 응답이 올바르게 재사용되도록 함
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 실제 요청 처리
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Vary", "Origin")
		}

		next.ServeHTTP(w, r)
	})
}

// isAllowedOrigin origin이 허용된 도메인인지 확인
func (c *CORS) isAllowedOrigin(origin string) bool {
	if len(c.allowedOrigins) == 0 {
		return true
	}

	for _, allowedOrigin := range c.allowedOrigins {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
		// 와일드카드 서브도메인 지원 (예: *.example.com)
		if strings.HasPrefix(allowedOrigin, "*.") {
			domain := strings.TrimPrefix(allowedOrigin, "*")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}

	return false
}
