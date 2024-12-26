package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/models"
	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/services"
	"github.com/kihyun1998/prisma-market/prisma-user-service/pkg/utils"
)

type UserHandler struct {
	userService *services.UserService
	jwtSecret   string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewUserHandler(userService *services.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

// CreateProfile 새로운 사용자 프로필 생성
func (h *UserHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// AuthID 파싱
	authID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		h.sendError(w, "Invalid auth ID", http.StatusBadRequest)
		return
	}

	var req models.CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.CreateProfile(r.Context(), authID, claims.Email, &req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile created successfully",
	})
}

// GetProfile 사용자 프로필 조회
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		h.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// GetProfileByUsername username으로 프로필 조회
func (h *UserHandler) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	profile, err := h.userService.GetProfileByUsername(r.Context(), username)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile 프로필 업데이트
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		h.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 프로필 접근 권한 확인
	profile, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		h.sendError(w, "Profile not found", http.StatusNotFound)
		return
	}

	if profile.AuthID.Hex() != claims.UserID {
		h.sendError(w, "Unauthorized to modify this profile", http.StatusForbidden)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateProfile(r.Context(), userID, &req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
	})
}

// DeleteProfile 프로필 삭제 (soft delete)
func (h *UserHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		h.sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 프로필 접근 권한 확인
	profile, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		h.sendError(w, "Profile not found", http.StatusNotFound)
		return
	}

	if profile.AuthID.Hex() != claims.UserID && claims.Role != "admin" {
		h.sendError(w, "Unauthorized to delete this profile", http.StatusForbidden)
		return
	}

	if err := h.userService.DeleteProfile(r.Context(), userID); err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile deleted successfully",
	})
}

// SearchProfiles 사용자 검색
func (h *UserHandler) SearchProfiles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	query = strings.TrimSpace(query)

	if query == "" {
		h.sendError(w, "Search query is required", http.StatusBadRequest)
		return
	}

	profiles, err := h.userService.SearchProfiles(r.Context(), query)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles)
}

// sendError 에러 응답 전송 헬퍼 함수
func (h *UserHandler) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
