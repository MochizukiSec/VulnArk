package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vulnark/vulnark/utils"
	"golang.org/x/crypto/bcrypt"
)

// User 简化版用户模型
type User struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `gorm:"type:varchar(50);unique_index;not null"`
	Password  string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);unique_index;not null"`
	RealName  string `gorm:"type:varchar(50)"`
	Role      string `gorm:"type:varchar(20);not null;default:'viewer'"`
	Active    bool   `gorm:"default:true"`
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func main() {
	// 数据库配置
	dsn := "root:root123456@tcp(127.0.0.1:3306)/vulnark?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	// 管理员信息
	adminUsername := "admin"
	adminPassword := "admin123"
	adminEmail := "admin@vulnark.com"
	adminRealName := "系统管理员"
	adminRole := "admin"

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 检查用户是否存在
	var count int
	db.Model(&User{}).Where("username = ? AND deleted_at IS NULL", adminUsername).Count(&count)
	if count > 0 {
		// 如果存在则软删除
		fmt.Printf("用户 %s 已存在，执行软删除...\n", adminUsername)
		db.Model(&User{}).Where("username = ? AND deleted_at IS NULL", adminUsername).
			Update("deleted_at", time.Now())
	}

	// 创建新管理员用户
	now := time.Now()
	admin := User{
		Username:  adminUsername,
		Password:  string(hashedPassword),
		Email:     adminEmail,
		RealName:  adminRealName,
		Role:      adminRole,
		Active:    true,
		LastLogin: now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("创建管理员失败: %v", err)
	}

	fmt.Println("✅ 管理员创建成功！")
	fmt.Printf("用户名: %s\n", adminUsername)
	fmt.Printf("密码: %s\n", adminPassword)
	fmt.Printf("角色: %s\n", adminRole)

	// 查询创建的管理员信息
	var result User
	db.Where("username = ? AND deleted_at IS NULL", adminUsername).First(&result)
	fmt.Printf("管理员ID: %d, 创建时间: %s\n", result.ID, utils.FormatTimeCST(result.CreatedAt))
}
