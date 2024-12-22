package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleUser   UserRole = "USER"
	RoleAdmin  UserRole = "ADMIN"
	RoleSeller UserRole = "SELLER"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"-"`
	Name        string             `bson:"name" json:"name"`
	Role        UserRole           `bson:"role" json:"role"`
	Phone       string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Address     string             `bson:"address,omitempty" json:"address,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	LastLoginAt time.Time          `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
}
