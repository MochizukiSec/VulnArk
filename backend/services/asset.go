package services

import (
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"vuln-management/config"
	"vuln-management/models"
)

// AssetService 處理資產相關邏輯
type AssetService struct {
	assetsColl *mongo.Collection
	vulnsColl  *mongo.Collection
}

// NewAssetService 創建資產服務實例
func NewAssetService(db *mongo.Database) *AssetService {
	return &AssetService{
		assetsColl: config.GetCollection(config.AssetsCollection),
		vulnsColl:  config.GetCollection(config.VulnerabilitiesCollection),
	}
}

// CreateAsset 創建新資產
func (s *AssetService) CreateAsset(ctx context.Context, create models.AssetCreate, userID primitive.ObjectID) (*models.Asset, error) {
	now := time.Now()
	asset := models.Asset{
		ID:           primitive.NewObjectID(),
		Name:         create.Name,
		Description:  create.Description,
		Type:         create.Type,
		Status:       create.Status,
		IPAddress:    create.IPAddress,
		MACAddress:   create.MACAddress,
		Location:     create.Location,
		Owner:        create.Owner,
		Department:   create.Department,
		PurchaseDate: create.PurchaseDate,
		ExpiryDate:   create.ExpiryDate,
		OS:           create.OS,
		OSVersion:    create.OSVersion,
		Manufacturer: create.Manufacturer,
		Model:        create.Model,
		SerialNumber: create.SerialNumber,
		Tags:         create.Tags,
		CustomFields: create.CustomFields,
		CreatedBy:    userID,
		UpdatedBy:    userID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err := s.assetsColl.InsertOne(ctx, asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetAssetByID 根據ID獲取資產
func (s *AssetService) GetAssetByID(ctx context.Context, id string) (*models.Asset, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var asset models.Asset
	err = s.assetsColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetAssetResponse 獲取包含漏洞數量的資產響應
func (s *AssetService) GetAssetResponse(ctx context.Context, id string) (*models.AssetResponse, error) {
	asset, err := s.GetAssetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 獲取漏洞數量
	vulnCount := len(asset.Vulnerabilities)

	// 構造響應
	response := models.AssetResponse{
		ID:           asset.ID.Hex(),
		Name:         asset.Name,
		Description:  asset.Description,
		Type:         asset.Type,
		Status:       asset.Status,
		IPAddress:    asset.IPAddress,
		MACAddress:   asset.MACAddress,
		Location:     asset.Location,
		Owner:        asset.Owner,
		Department:   asset.Department,
		PurchaseDate: asset.PurchaseDate,
		ExpiryDate:   asset.ExpiryDate,
		OS:           asset.OS,
		OSVersion:    asset.OSVersion,
		Manufacturer: asset.Manufacturer,
		Model:        asset.Model,
		SerialNumber: asset.SerialNumber,
		VulnCount:    vulnCount,
		Tags:         asset.Tags,
		CustomFields: asset.CustomFields,
		CreatedAt:    asset.CreatedAt,
		UpdatedAt:    asset.UpdatedAt,
	}

	return &response, nil
}

// UpdateAsset 更新資產信息
func (s *AssetService) UpdateAsset(ctx context.Context, id string, update models.AssetUpdate, userID primitive.ObjectID) (*models.Asset, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// 構建更新文檔
	updateDoc := bson.M{"updated_at": time.Now(), "updated_by": userID}

	if update.Name != nil {
		updateDoc["name"] = *update.Name
	}
	if update.Description != nil {
		updateDoc["description"] = *update.Description
	}
	if update.Type != nil {
		updateDoc["type"] = *update.Type
	}
	if update.Status != nil {
		updateDoc["status"] = *update.Status
	}
	if update.IPAddress != nil {
		updateDoc["ip_address"] = *update.IPAddress
	}
	if update.MACAddress != nil {
		updateDoc["mac_address"] = *update.MACAddress
	}
	if update.Location != nil {
		updateDoc["location"] = *update.Location
	}
	if update.Owner != nil {
		updateDoc["owner"] = *update.Owner
	}
	if update.Department != nil {
		updateDoc["department"] = *update.Department
	}
	if update.PurchaseDate != nil {
		updateDoc["purchase_date"] = update.PurchaseDate
	}
	if update.ExpiryDate != nil {
		updateDoc["expiry_date"] = update.ExpiryDate
	}
	if update.OS != nil {
		updateDoc["os"] = *update.OS
	}
	if update.OSVersion != nil {
		updateDoc["os_version"] = *update.OSVersion
	}
	if update.Manufacturer != nil {
		updateDoc["manufacturer"] = *update.Manufacturer
	}
	if update.Model != nil {
		updateDoc["model"] = *update.Model
	}
	if update.SerialNumber != nil {
		updateDoc["serial_number"] = *update.SerialNumber
	}
	if update.Tags != nil {
		updateDoc["tags"] = update.Tags
	}
	if update.CustomFields != nil {
		updateDoc["custom_fields"] = update.CustomFields
	}

	// 執行更新
	result := s.assetsColl.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateDoc},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	// 解碼更新後的資產
	var updatedAsset models.Asset
	if err := result.Decode(&updatedAsset); err != nil {
		return nil, err
	}

	return &updatedAsset, nil
}

// DeleteAsset 刪除資產
func (s *AssetService) DeleteAsset(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.assetsColl.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// SearchAssets 搜索資產
func (s *AssetService) SearchAssets(ctx context.Context, params models.AssetSearchParams) (*models.AssetListResponse, error) {
	// 構建查詢條件
	filter := bson.M{}

	if params.Type != "" {
		filter["type"] = params.Type
	}
	if params.Status != "" {
		filter["status"] = params.Status
	}
	if params.Department != "" {
		filter["department"] = params.Department
	}
	if params.SearchTerm != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": params.SearchTerm, "$options": "i"}},
			{"description": bson.M{"$regex": params.SearchTerm, "$options": "i"}},
			{"ip_address": bson.M{"$regex": params.SearchTerm, "$options": "i"}},
			{"mac_address": bson.M{"$regex": params.SearchTerm, "$options": "i"}},
			{"serial_number": bson.M{"$regex": params.SearchTerm, "$options": "i"}},
		}
	}

	// 設置分頁參數
	page := params.Page
	if page < 1 {
		page = 1
	}
	perPage := params.PerPage
	if perPage < 1 {
		perPage = 20
	}

	// 設置排序
	sortBy := "created_at"
	if params.SortBy != "" {
		// 映射前端字段到數據庫字段
		fieldMapping := map[string]string{
			"name":       "name",
			"type":       "type",
			"status":     "status",
			"department": "department",
			"location":   "location",
			"createdAt":  "created_at",
			"updatedAt":  "updated_at",
		}
		if mapped, ok := fieldMapping[params.SortBy]; ok {
			sortBy = mapped
		}
	}

	sortOrder := -1 // 默認降序
	if params.SortOrder == "asc" {
		sortOrder = 1
	}

	// 執行查詢
	skip := (page - 1) * perPage
	findOptions := options.Find().
		SetSort(bson.M{sortBy: sortOrder}).
		SetSkip(int64(skip)).
		SetLimit(int64(perPage))

	cursor, err := s.assetsColl.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 解析結果
	var assets []models.Asset
	if err := cursor.All(ctx, &assets); err != nil {
		return nil, err
	}

	// 獲取總數
	total, err := s.assetsColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 構建分頁信息
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	pagination := models.Pagination{
		Total:    int(total),
		Page:     page,
		PerPage:  perPage,
		LastPage: lastPage,
	}

	// 轉換為響應格式
	assetResponses := make([]models.AssetResponse, len(assets))
	for i, asset := range assets {
		assetResponses[i] = models.AssetResponse{
			ID:           asset.ID.Hex(),
			Name:         asset.Name,
			Description:  asset.Description,
			Type:         asset.Type,
			Status:       asset.Status,
			IPAddress:    asset.IPAddress,
			MACAddress:   asset.MACAddress,
			Location:     asset.Location,
			Owner:        asset.Owner,
			Department:   asset.Department,
			PurchaseDate: asset.PurchaseDate,
			ExpiryDate:   asset.ExpiryDate,
			OS:           asset.OS,
			OSVersion:    asset.OSVersion,
			Manufacturer: asset.Manufacturer,
			Model:        asset.Model,
			SerialNumber: asset.SerialNumber,
			VulnCount:    len(asset.Vulnerabilities),
			Tags:         asset.Tags,
			CustomFields: asset.CustomFields,
			CreatedAt:    asset.CreatedAt,
			UpdatedAt:    asset.UpdatedAt,
		}
	}

	return &models.AssetListResponse{
		Assets:     assetResponses,
		Pagination: pagination,
	}, nil
}

// AddVulnerabilityToAsset 將漏洞關聯到資產
func (s *AssetService) AddVulnerabilityToAsset(ctx context.Context, assetID string, vulnID string) error {
	assetObjID, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return err
	}

	vulnObjID, err := primitive.ObjectIDFromHex(vulnID)
	if err != nil {
		return err
	}

	// 更新資產，添加漏洞ID到列表
	_, err = s.assetsColl.UpdateOne(
		ctx,
		bson.M{"_id": assetObjID},
		bson.M{
			"$addToSet": bson.M{"vulnerabilities": vulnObjID},
			"$set":      bson.M{"updated_at": time.Now()},
		},
	)

	return err
}

// RemoveVulnerabilityFromAsset 從資產中移除漏洞關聯
func (s *AssetService) RemoveVulnerabilityFromAsset(ctx context.Context, assetID string, vulnID string) error {
	assetObjID, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return err
	}

	vulnObjID, err := primitive.ObjectIDFromHex(vulnID)
	if err != nil {
		return err
	}

	// 更新資產，從列表中移除漏洞ID
	_, err = s.assetsColl.UpdateOne(
		ctx,
		bson.M{"_id": assetObjID},
		bson.M{
			"$pull": bson.M{"vulnerabilities": vulnObjID},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)

	return err
}

// GetAssetVulnerabilities 獲取資產關聯的所有漏洞
func (s *AssetService) GetAssetVulnerabilities(ctx context.Context, assetID string) ([]models.Vulnerability, error) {
	// 先獲取資產
	asset, err := s.GetAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	// 如果沒有關聯漏洞，返回空數組
	if len(asset.Vulnerabilities) == 0 {
		return []models.Vulnerability{}, nil
	}

	// 查詢所有關聯的漏洞
	cursor, err := s.vulnsColl.Find(
		ctx,
		bson.M{"_id": bson.M{"$in": asset.Vulnerabilities}},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 解析結果
	var vulnerabilities []models.Vulnerability
	if err := cursor.All(ctx, &vulnerabilities); err != nil {
		return nil, err
	}

	return vulnerabilities, nil
}

// AddAssetNote 給資產添加備註
func (s *AssetService) AddAssetNote(ctx context.Context, assetID string, content string, userID primitive.ObjectID) error {
	assetObjID, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return err
	}

	note := models.AssetNote{
		ID:        primitive.NewObjectID(),
		Content:   content,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	_, err = s.assetsColl.UpdateOne(
		ctx,
		bson.M{"_id": assetObjID},
		bson.M{
			"$push": bson.M{"notes": note},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)

	return err
}
