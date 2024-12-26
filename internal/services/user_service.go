package services

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/models"
	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo *mongodb.UserRepository
}

func NewUserService(repo *mongodb.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateProfile(ctx context.Context, authID primitive.ObjectID, email string, req *models.CreateProfileRequest) error {
	// 입력값 검증
	if err := validateCreateRequest(req); err != nil {
		return err
	}

	// username 중복 체크
	existingProfile, err := s.repo.GetProfileByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if existingProfile != nil {
		return errors.New("username already exists")
	}

	// 프로필 생성
	profile := &models.UserProfile{
		AuthID:      authID,
		Email:       email,
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
	}

	return s.repo.CreateProfile(ctx, profile)
}

func (s *UserService) GetProfile(ctx context.Context, userID primitive.ObjectID) (*models.ProfileResponse, error) {
	profile, err := s.repo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	return &models.ProfileResponse{
		UserProfile: *profile,
	}, nil
}

func (s *UserService) GetProfileByUsername(ctx context.Context, username string) (*models.ProfileResponse, error) {
	profile, err := s.repo.GetProfileByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	return &models.ProfileResponse{
		UserProfile: *profile,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID primitive.ObjectID, req *models.UpdateProfileRequest) error {
	// 프로필 존재 확인
	profile, err := s.repo.GetProfileByID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	// 업데이트할 필드 수집
	update := bson.M{}

	if req.Username != nil {
		if err := validateUsername(*req.Username); err != nil {
			return err
		}
		// username 중복 체크
		existingProfile, err := s.repo.GetProfileByUsername(ctx, *req.Username)
		if err != nil {
			return err
		}
		if existingProfile != nil && existingProfile.ID != userID {
			return errors.New("username already exists")
		}
		update["username"] = *req.Username
	}

	if req.FirstName != nil {
		if err := validateName(*req.FirstName); err != nil {
			return err
		}
		update["first_name"] = *req.FirstName
	}

	if req.LastName != nil {
		if err := validateName(*req.LastName); err != nil {
			return err
		}
		update["last_name"] = *req.LastName
	}

	if req.PhoneNumber != nil {
		if err := validatePhoneNumber(*req.PhoneNumber); err != nil {
			return err
		}
		update["phone_number"] = *req.PhoneNumber
	}

	if req.Address != nil {
		if err := validateAddress(req.Address); err != nil {
			return err
		}
		update["address"] = req.Address
	}

	if len(update) == 0 {
		return nil // 업데이트할 내용이 없음
	}

	return s.repo.UpdateProfile(ctx, userID, update)
}

func (s *UserService) DeleteProfile(ctx context.Context, userID primitive.ObjectID) error {
	return s.repo.DeleteProfile(ctx, userID)
}

func (s *UserService) SearchProfiles(ctx context.Context, query string) ([]*models.ProfileResponse, error) {
	if len(strings.TrimSpace(query)) < 2 {
		return nil, errors.New("search query must be at least 2 characters")
	}

	profiles, err := s.repo.SearchProfiles(ctx, query, 20) // 최대 20개 결과
	if err != nil {
		return nil, err
	}

	responses := make([]*models.ProfileResponse, len(profiles))
	for i, profile := range profiles {
		responses[i] = &models.ProfileResponse{
			UserProfile: *profile,
		}
	}

	return responses, nil
}

// Validation helpers
func validateCreateRequest(req *models.CreateProfileRequest) error {
	if err := validateUsername(req.Username); err != nil {
		return err
	}
	if err := validateName(req.FirstName); err != nil {
		return err
	}
	if err := validateName(req.LastName); err != nil {
		return err
	}
	if err := validatePhoneNumber(req.PhoneNumber); err != nil {
		return err
	}
	if err := validateAddress(&req.Address); err != nil {
		return err
	}
	return nil
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}
	if match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username); !match {
		return errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}
	return nil
}

func validateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 1 || len(name) > 50 {
		return errors.New("name must be between 1 and 50 characters")
	}
	if match, _ := regexp.MatchString("^[a-zA-Z\\s-]+$", name); !match {
		return errors.New("name can only contain letters, spaces, and hyphens")
	}
	return nil
}

func validatePhoneNumber(phone string) error {
	phone = strings.TrimSpace(phone)
	if match, _ := regexp.MatchString("^\\+?[0-9]{10,15}$", phone); !match {
		return errors.New("invalid phone number format")
	}
	return nil
}

func validateAddress(addr *models.Address) error {
	if addr == nil {
		return errors.New("address is required")
	}
	if strings.TrimSpace(addr.Street) == "" {
		return errors.New("street is required")
	}
	if strings.TrimSpace(addr.City) == "" {
		return errors.New("city is required")
	}
	if strings.TrimSpace(addr.Country) == "" {
		return errors.New("country is required")
	}
	return nil
}
