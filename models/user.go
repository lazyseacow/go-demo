package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // 密码不返回给前端
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	Status    int            `gorm:"type:tinyint;default:1;comment:状态 1正常 0禁用" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 这里可以添加创建前的逻辑，比如密码加密
	return nil
}
