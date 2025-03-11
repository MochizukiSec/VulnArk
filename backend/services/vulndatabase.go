package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"vuln-management/config"
	"vuln-management/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VulnDatabaseService 提供漏洞库相关的服务
type VulnDatabaseService struct {
	collection *mongo.Collection
}

// NewVulnDatabaseService 创建新的漏洞库服务实例
func NewVulnDatabaseService(db *mongo.Database) *VulnDatabaseService {
	return &VulnDatabaseService{
		collection: db.Collection(config.VulnDatabaseCollection),
	}
}

// toDocument 将漏洞对象转换为BSON文档
func toDocument(entry models.VulnDBEntry) bson.M {
	// 处理ID转换
	var objID primitive.ObjectID
	if entry.ID != "" {
		var err error
		objID, err = primitive.ObjectIDFromHex(entry.ID)
		if err != nil {
			// 如果ID格式不正确，创建新ID
			objID = primitive.NewObjectID()
		}
	} else {
		objID = primitive.NewObjectID()
	}

	// 确保时间戳存在
	now := time.Now()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	if entry.UpdatedAt.IsZero() {
		entry.UpdatedAt = now
	}

	return bson.M{
		"_id":              objID,
		"cveId":            entry.CveId,
		"title":            entry.Title,
		"description":      entry.Description,
		"cvss":             entry.Cvss,
		"severity":         entry.Severity,
		"affectedSystems":  entry.AffectedSystems,
		"solution":         entry.Solution,
		"publishedDate":    entry.PublishedDate,
		"lastModifiedDate": entry.LastModifiedDate,
		"references":       entry.References,
		"tags":             entry.Tags,
		"createdAt":        entry.CreatedAt,
		"updatedAt":        entry.UpdatedAt,
	}
}

// fromDocument 将BSON文档转换为漏洞对象
func fromDocument(doc bson.M) models.VulnDBEntry {
	id := ""
	if objID, ok := doc["_id"].(primitive.ObjectID); ok {
		id = objID.Hex()
	}

	var references []string
	if refsInterface, ok := doc["references"].(primitive.A); ok {
		for _, ref := range refsInterface {
			if refStr, ok := ref.(string); ok {
				references = append(references, refStr)
			}
		}
	}

	var tags []string
	if tagsInterface, ok := doc["tags"].(primitive.A); ok {
		for _, tag := range tagsInterface {
			if tagStr, ok := tag.(string); ok {
				tags = append(tags, tagStr)
			}
		}
	}

	createdAt := time.Now()
	if timeValue, ok := doc["createdAt"].(primitive.DateTime); ok {
		createdAt = timeValue.Time()
	}

	updatedAt := time.Now()
	if timeValue, ok := doc["updatedAt"].(primitive.DateTime); ok {
		updatedAt = timeValue.Time()
	}

	return models.VulnDBEntry{
		ID:               id,
		CveId:            getStringValue(doc, "cveId"),
		Title:            getStringValue(doc, "title"),
		Description:      getStringValue(doc, "description"),
		Cvss:             getFloat64Value(doc, "cvss"),
		Severity:         getStringValue(doc, "severity"),
		AffectedSystems:  getStringValue(doc, "affectedSystems"),
		Solution:         getStringValue(doc, "solution"),
		PublishedDate:    getStringValue(doc, "publishedDate"),
		LastModifiedDate: getStringValue(doc, "lastModifiedDate"),
		References:       references,
		Tags:             tags,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

// getStringValue 安全地从BSON文档获取字符串值
func getStringValue(doc bson.M, key string) string {
	if value, ok := doc[key].(string); ok {
		return value
	}
	return ""
}

// getFloat64Value 安全地从BSON文档获取浮点数值
func getFloat64Value(doc bson.M, key string) float64 {
	if value, ok := doc[key].(float64); ok {
		return value
	}
	return 0
}

// CreateVulnerability 创建新漏洞
func (s *VulnDatabaseService) CreateVulnerability(vulnEntry models.VulnDBEntry) (models.VulnDBEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置创建时间和更新时间
	now := time.Now()
	vulnEntry.CreatedAt = now
	vulnEntry.UpdatedAt = now

	// 转换为BSON文档
	doc := toDocument(vulnEntry)

	// 插入文档
	result, err := s.collection.InsertOne(ctx, doc)
	if err != nil {
		return models.VulnDBEntry{}, err
	}

	// 获取插入的ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		vulnEntry.ID = oid.Hex()
	}

	return vulnEntry, nil
}

// GetVulnerabilityByCveID 根据CVE ID获取漏洞
func (s *VulnDatabaseService) GetVulnerabilityByCveID(cveID string) (models.VulnDBEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建查询条件
	filter := bson.M{"cveId": cveID}

	// 执行查询
	var result bson.M
	err := s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.VulnDBEntry{}, fmt.Errorf("未找到CVE ID为 %s 的漏洞", cveID)
		}
		return models.VulnDBEntry{}, err
	}

	return fromDocument(result), nil
}

// GetVulnerabilityByID 根据ID获取漏洞
func (s *VulnDatabaseService) GetVulnerabilityByID(id string) (models.VulnDBEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 转换ID字符串为ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.VulnDBEntry{}, fmt.Errorf("无效的漏洞ID: %s", id)
	}

	// 创建查询条件
	filter := bson.M{"_id": objID}

	// 执行查询
	var result bson.M
	err = s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.VulnDBEntry{}, fmt.Errorf("未找到ID为 %s 的漏洞", id)
		}
		return models.VulnDBEntry{}, err
	}

	return fromDocument(result), nil
}

// SearchVulnerabilities 搜索漏洞库
func (s *VulnDatabaseService) SearchVulnerabilities(query string, severity string, yearStr string, sortBy string, sortOrder string, page int, perPage int) ([]models.VulnDBEntry, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 创建查询条件
	filter := bson.M{}

	// 关键词搜索
	if query != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"cveId": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	// 按严重程度筛选
	if severity != "" {
		filter["severity"] = severity
	}

	// 按年份筛选
	if yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			yearPrefix := fmt.Sprintf("%d-", year)
			filter["publishedDate"] = bson.M{"$regex": "^" + yearPrefix}
		}
	}

	// 设置排序
	sortOptions := bson.D{}

	// 设置默认排序
	if sortBy == "" {
		sortBy = "publishedDate"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// 将前端排序字段映射到MongoDB字段
	var mongoField string
	switch sortBy {
	case "publishedDate":
		mongoField = "publishedDate"
	case "lastModifiedDate":
		mongoField = "lastModifiedDate"
	case "cveId":
		mongoField = "cveId"
	case "cvss":
		mongoField = "cvss"
	default:
		mongoField = "publishedDate"
	}

	// 设置排序方向
	sortValue := 1
	if sortOrder == "desc" {
		sortValue = -1
	}

	sortOptions = append(sortOptions, bson.E{Key: mongoField, Value: sortValue})

	// 计算总记录数
	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 设置分页
	skip := (page - 1) * perPage
	if skip < 0 {
		skip = 0
	}

	// 设置查询选项
	findOptions := options.Find().
		SetSort(sortOptions).
		SetSkip(int64(skip)).
		SetLimit(int64(perPage))

	// 执行查询
	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 解析结果
	var results []models.VulnDBEntry
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		results = append(results, fromDocument(doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return results, int(total), nil
}

// UpdateVulnerability 更新漏洞
func (s *VulnDatabaseService) UpdateVulnerability(cveID string, vulnEntry models.VulnDBEntry) (models.VulnDBEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置更新时间
	vulnEntry.UpdatedAt = time.Now()

	// 转换为BSON文档，但不包含ID和创建时间
	updateDoc := toDocument(vulnEntry)
	delete(updateDoc, "_id")       // ID不应被更新
	delete(updateDoc, "createdAt") // 创建时间不应被更新

	// 创建更新操作
	update := bson.M{
		"$set": updateDoc,
	}

	// 执行更新
	filter := bson.M{"cveId": cveID}
	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.VulnDBEntry{}, err
	}

	if result.MatchedCount == 0 {
		return models.VulnDBEntry{}, fmt.Errorf("未找到CVE ID为 %s 的漏洞", cveID)
	}

	// 返回更新后的文档
	return s.GetVulnerabilityByCveID(cveID)
}

// DeleteVulnerability 删除漏洞
func (s *VulnDatabaseService) DeleteVulnerability(cveID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建删除条件
	filter := bson.M{"cveId": cveID}

	// 执行删除
	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("未找到CVE ID为 %s 的漏洞", cveID)
	}

	return nil
}

// ImportInitialData 导入初始数据
func (s *VulnDatabaseService) ImportInitialData(entries []models.VulnDBEntry) error {
	// 检查集合是否为空
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	// 如果集合不为空，不导入初始数据
	if count > 0 {
		return nil
	}

	// 准备批量插入的文档
	var documents []interface{}
	for _, entry := range entries {
		if entry.CreatedAt.IsZero() {
			entry.CreatedAt = time.Now()
		}
		if entry.UpdatedAt.IsZero() {
			entry.UpdatedAt = time.Now()
		}
		documents = append(documents, toDocument(entry))
	}

	// 批量插入
	_, err = s.collection.InsertMany(ctx, documents)
	return err
}
