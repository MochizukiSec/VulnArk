package services

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// ReportService 报告相关业务逻辑
type ReportService struct {
	db             *mongo.Database
	collectionName string
}

// NewReportService 创建报告服务
func NewReportService(db *mongo.Database) *ReportService {
	return &ReportService{
		db:             db,
		collectionName: "reports",
	}
}

// 添加其他报告服务方法
// ...
