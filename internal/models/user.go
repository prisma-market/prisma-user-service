package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserProfile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AuthID      primitive.ObjectID `bson:"auth_id" json:"auth_id"` // Auth Service의 사용자 ID
	Email       string             `bson:"email" json:"email"`     // Auth Service와 동기화
	Username    string             `bson:"username" json:"username"`
	FirstName   string             `bson:"first_name" json:"first_name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Address     Address            `bson:"address" json:"address"`
	Avatar      string             `bson:"avatar" json:"avatar"` // 프로필 이미지 URL
	Status      string             `bson:"status" json:"status"` // active, inactive, suspended
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Address struct {
	Street     string `bson:"street" json:"street"`
	City       string `bson:"city" json:"city"`
	State      string `bson:"state" json:"state"`
	PostalCode string `bson:"postal_code" json:"postal_code"`
	Country    string `bson:"country" json:"country"`
}

// API 요청/응답 구조체
type CreateProfileRequest struct {
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber string  `json:"phone_number"`
	Address     Address `json:"address"`
}

type UpdateProfileRequest struct {
	Username    *string  `json:"username,omitempty"`
	FirstName   *string  `json:"first_name,omitempty"`
	LastName    *string  `json:"last_name,omitempty"`
	PhoneNumber *string  `json:"phone_number,omitempty"`
	Address     *Address `json:"address,omitempty"`
}

type ProfileResponse struct {
	UserProfile
	// 추가 필드가 필요한 경우 여기에 정의
}
