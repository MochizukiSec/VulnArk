package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// 用户角色
const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleViewer = "viewer"
)

// 用户状态
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusDisabled = "disabled"
)

// User 表示用户数据模型
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username"`
	Email          string             `bson:"email" json:"email"`
	PasswordHash   string             `bson:"password_hash" json:"-"`
	FirstName      string             `bson:"first_name" json:"firstName"`
	LastName       string             `bson:"last_name" json:"lastName"`
	Role           string             `bson:"role" json:"role"`
	Department     string             `bson:"department" json:"department"`
	ProfilePicture string             `bson:"profile_picture" json:"profilePicture"`
	Status         string             `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updatedAt"`
	LastLogin      time.Time          `bson:"last_login" json:"lastLogin"`
}

// UserCredentials 用于认证的用户凭证
type UserCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserRegistration 用于注册的用户数据
type UserRegistration struct {
	Username   string `json:"username" binding:"required,min=3,max=20"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Department string `json:"department"`
	Role       string `json:"role"`
	Status     string `json:"status"`
}

// UserResponse 是返回给客户端的用户信息
type UserResponse struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	Department     string    `json:"department"`
	ProfilePicture string    `json:"profile_picture"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	LastLogin      time.Time `json:"last_login"`
}

// SetPassword 设置用户密码
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword 验证用户密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// ToResponse 将用户模型转换为响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:             u.ID.Hex(),
		Username:       u.Username,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Name:           u.FirstName + " " + u.LastName,
		Role:           u.Role,
		Department:     u.Department,
		ProfilePicture: u.ProfilePicture,
		Status:         u.Status,
		CreatedAt:      u.CreatedAt,
		LastLogin:      u.LastLogin,
	}
}
