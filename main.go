package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prisma-market/prisma-user-service/handler"
	"github.com/prisma-market/prisma-user-service/repository"
	"github.com/prisma-market/prisma-user-service/service"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	userHandler *handler.UserHandler
)

func connectDB(mongoURI string) (*mongo.Client, error) {
	maxRetries := 5
	var client *mongo.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		// MongoDB 연결
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err == nil {
			// 연결 테스트
			if err = client.Ping(ctx, nil); err == nil {
				log.Printf("Successfully connected to MongoDB on attempt %d", i+1)
				return client, nil
			}
		}

		log.Printf("Failed to connect to MongoDB (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to MongoDB after %d attempts", maxRetries)
}

func init() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// MongoDB 연결 시도
	var err error
	mongoClient, err = connectDB(mongoURI)
	if err != nil {
		log.Fatal("Could not establish MongoDB connection:", err)
	}

	// 데이터베이스 및 의존성 초기화
	db := mongoClient.Database("shopping_mall")
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler = handler.NewUserHandler(userService)

	log.Println("Successfully initialized MongoDB connection and services")
}

func main() {
	// MongoDB 연결 종료
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	// 라우터 설정
	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	// 미들웨어 설정
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	// 서버 설정
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 서버 시작
	go func() {
		log.Printf("Server is running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
