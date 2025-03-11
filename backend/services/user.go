package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService 用户相关业务逻辑
type UserService struct {
	db             *mongo.Database
	collectionName string
}

// NewUserService 创建用户服务
func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		db:             db,
		collectionName: "users",
	}
}

// 添加其他用户服务方法
// ...
