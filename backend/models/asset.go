package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Pagination 分頁結構
type Pagination struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PerPage  int `json:"perPage"`
	LastPage int `json:"lastPage"`
}

// 資產類型常量
const (
	AssetTypeServer      = "server"
	AssetTypeWorkstation = "workstation"
	AssetTypeNetwork     = "network"
	AssetTypeApplication = "application"
	AssetTypeDatabase    = "database"
	AssetTypeCloud       = "cloud"
	AssetTypeContainer   = "container"
	AssetTypeIoT         = "iot"
	AssetTypeOther       = "other"
)

// 資產狀態常量
const (
	AssetStatusActive         = "active"
	AssetStatusMaintenance    = "maintenance"
	AssetStatusDecommissioned = "decommissioned"
	AssetStatusReserved       = "reserved"
	AssetStatusIssue          = "issue"
)

// Asset 表示資產數據模型
type Asset struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name            string                 `bson:"name" json:"name"`
	Description     string                 `bson:"description" json:"description"`
	Type            string                 `bson:"type" json:"type"`
	Status          string                 `bson:"status" json:"status"`
	IPAddress       string                 `bson:"ip_address" json:"ipAddress"`
	MACAddress      string                 `bson:"mac_address" json:"macAddress"`
	Location        string                 `bson:"location" json:"location"`
	Owner           string                 `bson:"owner" json:"owner"`
	OwnerID         primitive.ObjectID     `bson:"owner_id,omitempty" json:"ownerId"`
	Department      string                 `bson:"department" json:"department"`
	PurchaseDate    *time.Time             `bson:"purchase_date,omitempty" json:"purchaseDate"`
	ExpiryDate      *time.Time             `bson:"expiry_date,omitempty" json:"expiryDate"`
	OS              string                 `bson:"os" json:"os"`
	OSVersion       string                 `bson:"os_version" json:"osVersion"`
	Manufacturer    string                 `bson:"manufacturer" json:"manufacturer"`
	Model           string                 `bson:"model" json:"model"`
	SerialNumber    string                 `bson:"serial_number" json:"serialNumber"`
	Vulnerabilities []primitive.ObjectID   `bson:"vulnerabilities" json:"vulnerabilities"`
	Tags            []string               `bson:"tags" json:"tags"`
	Notes           []AssetNote            `bson:"notes" json:"notes"`
	CustomFields    map[string]interface{} `bson:"custom_fields" json:"customFields"`
	CreatedBy       primitive.ObjectID     `bson:"created_by" json:"createdBy"`
	UpdatedBy       primitive.ObjectID     `bson:"updated_by" json:"updatedBy"`
	CreatedAt       time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updatedAt"`
}

// AssetNote 表示關於資產的備註
type AssetNote struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Content   string             `bson:"content" json:"content"`
	CreatedBy primitive.ObjectID `bson:"created_by" json:"createdBy"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// AssetCreate 用於創建資產的請求結構
type AssetCreate struct {
	Name         string                 `json:"name" binding:"required"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type" binding:"required"`
	Status       string                 `json:"status" binding:"required"`
	IPAddress    string                 `json:"ipAddress"`
	MACAddress   string                 `json:"macAddress"`
	Location     string                 `json:"location"`
	Owner        string                 `json:"owner"`
	Department   string                 `json:"department"`
	PurchaseDate *time.Time             `json:"purchaseDate"`
	ExpiryDate   *time.Time             `json:"expiryDate"`
	OS           string                 `json:"os"`
	OSVersion    string                 `json:"osVersion"`
	Manufacturer string                 `json:"manufacturer"`
	Model        string                 `json:"model"`
	SerialNumber string                 `json:"serialNumber"`
	Tags         []string               `json:"tags"`
	CustomFields map[string]interface{} `json:"customFields"`
}

// AssetUpdate 用於更新資產的請求結構
type AssetUpdate struct {
	Name         *string                `json:"name"`
	Description  *string                `json:"description"`
	Type         *string                `json:"type"`
	Status       *string                `json:"status"`
	IPAddress    *string                `json:"ipAddress"`
	MACAddress   *string                `json:"macAddress"`
	Location     *string                `json:"location"`
	Owner        *string                `json:"owner"`
	Department   *string                `json:"department"`
	PurchaseDate *time.Time             `json:"purchaseDate"`
	ExpiryDate   *time.Time             `json:"expiryDate"`
	OS           *string                `json:"os"`
	OSVersion    *string                `json:"osVersion"`
	Manufacturer *string                `json:"manufacturer"`
	Model        *string                `json:"model"`
	SerialNumber *string                `json:"serialNumber"`
	Tags         []string               `json:"tags"`
	CustomFields map[string]interface{} `json:"customFields"`
}

// AssetSearchParams 用於搜索資產的參數
type AssetSearchParams struct {
	Type       string `form:"type"`
	Status     string `form:"status"`
	Department string `form:"department"`
	SearchTerm string `form:"q"`
	SortBy     string `form:"sortBy"`
	SortOrder  string `form:"sortOrder"`
	Page       int    `form:"page,default=1"`
	PerPage    int    `form:"perPage,default=20"`
}

// AddNote 給資產添加一條備註
func (a *Asset) AddNote(content string, userID primitive.ObjectID) {
	note := AssetNote{
		ID:        primitive.NewObjectID(),
		Content:   content,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	a.Notes = append(a.Notes, note)
	a.UpdatedAt = time.Now()
}

// AssetResponse 用於返回資產資訊的響應結構
type AssetResponse struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	Status       string                 `json:"status"`
	IPAddress    string                 `json:"ipAddress"`
	MACAddress   string                 `json:"macAddress"`
	Location     string                 `json:"location"`
	Owner        string                 `json:"owner"`
	Department   string                 `json:"department"`
	PurchaseDate *time.Time             `json:"purchaseDate"`
	ExpiryDate   *time.Time             `json:"expiryDate"`
	OS           string                 `json:"os"`
	OSVersion    string                 `json:"osVersion"`
	Manufacturer string                 `json:"manufacturer"`
	Model        string                 `json:"model"`
	SerialNumber string                 `json:"serialNumber"`
	VulnCount    int                    `json:"vulnCount"`
	Tags         []string               `json:"tags"`
	CustomFields map[string]interface{} `json:"customFields"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// AssetListResponse 用於資產列表分頁響應
type AssetListResponse struct {
	Assets     []AssetResponse `json:"assets"`
	Pagination Pagination      `json:"pagination"`
}
