package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/config"
	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/handlers"
	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/repository/mongodb"
	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/services"
	"github.com/kihyun1998/prisma-market/prisma-user-service/pkg/middleware"
)

func main() {
	// 설정 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// MongoDB 리포지토리 초기화
	userRepo, err := mongodb.NewUserRepository(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	// 서비스 초기화
	userService := services.NewUserService(userRepo)

	// 핸들러 초기화
	userHandler := handlers.NewUserHandler(userService, cfg.JWTSecret)

	// 라우터 설정
	r := mux.NewRouter()

	// JWT 미들웨어 설정
	auth := middleware.NewJWTMiddleware(cfg.JWTSecret)

	// Public endpoints (인증 불필요)
	publicRouter := r.PathPrefix("/api/v1/public").Subrouter()
	publicRouter.HandleFunc("/users/search", userHandler.SearchProfiles).Methods("GET")
	publicRouter.HandleFunc("/users/username/{username}", userHandler.GetProfileByUsername).Methods("GET")
	publicRouter.HandleFunc("/users/{id}", userHandler.GetProfile).Methods("GET")

	// Protected endpoints (사용자 인증 필요)
	protectedRouter := r.PathPrefix("/api/v1").Subrouter()
	protectedRouter.Use(auth.ValidateJWT)
	protectedRouter.HandleFunc("/users", userHandler.CreateProfile).Methods("POST")
	protectedRouter.HandleFunc("/users/{id}", userHandler.UpdateProfile).Methods("PUT")
	protectedRouter.HandleFunc("/users/{id}", userHandler.DeleteProfile).Methods("DELETE")

	// Admin endpoints (관리자 권한 필요)
	adminRouter := r.PathPrefix("/api/v1/admin").Subrouter()
	adminRouter.Use(auth.ValidateJWT)
	adminRouter.Use(auth.RequireRole("admin"))
	// 관리자용 API 엔드포인트 추가

	// CORS 미들웨어 추가
	corsMiddleware := middleware.NewCORS()
	if len(cfg.AllowedOrigins) > 0 {
		corsMiddleware.WithOrigins(cfg.AllowedOrigins...)
	}
	r.Use(corsMiddleware.Handler)

	// 서버 시작
	log.Printf("Starting User Service on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
