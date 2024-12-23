package service

import (
	"context"
	"errors"
	"log"

	"github.com/prisma-market/prisma-user-service/models"
	"github.com/prisma-market/prisma-user-service/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type CreateUserInput struct {
	Email    string
	Password string
	Name     string
	Role     models.UserRole
	Phone    string
	Address  string
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	// Check if email already exists
	if _, err := s.repo.GetByEmail(ctx, input.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Name:     input.Name,
		Role:     input.Role,
		Phone:    input.Phone,
		Address:  input.Address,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	if phone, ok := updates["phone"].(string); ok {
		user.Phone = phone
	}
	if address, ok := updates["address"].(string); ok {
		user.Address = address
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		user.IsActive = isActive
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) ValidateCredentials(ctx context.Context, email, password string) (*models.User, error) {
	/// 1. 이메일로 사용자 조회
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. 비활성 사용자 체크
	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	// 3. 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 4. 마지막 로그인 시간 업데이트
	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		log.Printf("Failed to update last login: %v", err)
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int64) ([]models.User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.List(ctx, page, limit)
}

// 사용자 역할 변경
func (s *UserService) ChangeUserRole(ctx context.Context, userID primitive.ObjectID, newRole models.UserRole) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Role = newRole
	return s.repo.Update(ctx, user)
}
