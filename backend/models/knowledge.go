package models

import (
	"time"
)

// 知识类型
type KnowledgeType string

const (
	KnowledgeTypeTutorial   KnowledgeType = "tutorial"   // 教程
	KnowledgeTypeGuide      KnowledgeType = "guide"      // 指南
	KnowledgeTypeReference  KnowledgeType = "reference"  // 参考资料
	KnowledgeTypeMitigation KnowledgeType = "mitigation" // 缓解策略
	KnowledgeTypeProc       KnowledgeType = "procedure"  // 操作流程
	KnowledgeTypeOther      KnowledgeType = "other"      // 其他
)

// Knowledge 知识库条目模型
type Knowledge struct {
	ID      uint          `json:"id" gorm:"primary_key"`
	Title   string        `json:"title" gorm:"type:varchar(255);not null"`
	Content string        `json:"content" gorm:"type:text;not null"`
	Type    KnowledgeType `json:"type" gorm:"type:varchar(50);not null"`

	Tags       string `json:"tags" gorm:"type:varchar(255)"`
	Categories string `json:"categories" gorm:"type:varchar(255)"`

	RelatedVulnTypes string `json:"related_vuln_types" gorm:"type:varchar(255)"`

	Author   string `json:"author" gorm:"type:varchar(100)"`
	AuthorID uint   `json:"author_id" gorm:"index"`

	ViewCount int `json:"view_count" gorm:"default:0"`

	Attachments string `json:"attachments" gorm:"type:text"`
	References  string `json:"references" gorm:"type:text"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Knowledge) TableName() string {
	return "knowledge"
}

// IncrementViewCount 增加浏览次数
func (k *Knowledge) IncrementViewCount() {
	k.ViewCount++
}

// IsTutorial 判断是否为教程
func (k *Knowledge) IsTutorial() bool {
	return k.Type == KnowledgeTypeTutorial
}

// IsGuide 判断是否为指南
func (k *Knowledge) IsGuide() bool {
	return k.Type == KnowledgeTypeGuide
}

// IsMitigation 判断是否为缓解策略
func (k *Knowledge) IsMitigation() bool {
	return k.Type == KnowledgeTypeMitigation
}
