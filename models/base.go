package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PageRequest 分页请求参数
type PageRequest struct {
	Page     int `form:"page" json:"page" binding:"min=1"`                   // 页码
	PageSize int `form:"page_size" json:"page_size" binding:"min=1,max=100"` // 每页数量
}

// PageResponse 分页响应
type PageResponse struct {
	Total    int64 `json:"total"`     // 总数
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页数量
	List     any   `json:"list"`      // 数据列表
}

// DefaultPageRequest 默认分页参数
func DefaultPageRequest() PageRequest {
	return PageRequest{
		Page:     1,
		PageSize: 10,
	}
}
