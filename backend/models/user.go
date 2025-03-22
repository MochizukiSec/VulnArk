package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// 角色类型
type Role string

const (
	RoleAdmin    Role = "admin"    // 管理员
	RoleManager  Role = "manager"  // 经理
	RoleAuditor  Role = "auditor"  // 审计员
	RoleOperator Role = "operator" // 操作员
	RoleViewer   Role = "viewer"   // 浏览者
)

// User 用户模型
type User struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	Username  string     `json:"username" gorm:"type:varchar(50);unique_index;not null"`
	Email     string     `json:"email" gorm:"type:varchar(100);unique_index;not null"`
	Password  string     `json:"-" gorm:"type:varchar(100);not null"`
	RealName  string     `json:"real_name" gorm:"type:varchar(50)"`
	Phone     string     `json:"phone" gorm:"type:varchar(20)"`
	Role      Role       `json:"role" gorm:"type:varchar(20);not null;default:'viewer'"`
	Avatar    string     `json:"avatar" gorm:"type:varchar(255)"`
	LastLogin time.Time  `json:"last_login"`
	Active    bool       `json:"active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	// 检查密码字段是否在更新字段中
	if _, ok := scope.Get("gorm:update_column"); ok {
		// 当明确指定了更新字段时，检查是否包含密码字段
		for _, field := range scope.Fields() {
			if field.Name == "Password" && field.IsNormal && !field.IsIgnored {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
				if err != nil {
					return err
				}
				u.Password = string(hashedPassword)
				break
			}
		}
	} else {
		// 没有明确指定更新字段，如果密码不为空则更新
		if u.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			u.Password = string(hashedPassword)
		}
	}
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	// 如果密码哈希为空，直接返回false
	if u.Password == "" {
		fmt.Println("密码验证失败: 存储的密码哈希为空")
		return false
	}

	// 调试日志
	fmt.Printf("密码验证 - 用户ID: %d, 用户名: %s, 密码长度: %d\n",
		u.ID, u.Username, len(u.Password))

	// 使用bcrypt的CompareHashAndPassword函数比较密码
	// 该函数内部实现了恒定时间比较，可以有效防止时序攻击
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	// 打印验证结果
	if err != nil {
		fmt.Printf("密码验证失败: %v\n", err)
		return false
	}

	fmt.Println("密码验证成功")
	return true
}

// IsAdmin 判断是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// CanManage 判断是否有管理权限
func (u *User) CanManage() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}

// CanEdit 判断是否有编辑权限
func (u *User) CanEdit() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager || u.Role == RoleOperator
}

// CanAudit 判断是否有审计权限
func (u *User) CanAudit() bool {
	return u.Role == RoleAdmin || u.Role == RoleAuditor
}
