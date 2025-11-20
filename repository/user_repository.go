package repository

import (
	"github.com/demo/models"
	"gorm.io/gorm"
)

// UserRepository 用户仓库
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// FindByID 根据 ID 查找用户
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// Update 更新用户
func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

// UpdateFields 更新用户指定字段
func (r *UserRepository) UpdateFields(id int64, fields map[string]any) error {
	return r.DB.Model(&models.User{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(id int64) error {
	return r.DB.Delete(&models.User{}, id).Error
}

// List 获取用户列表（分页）
func (r *UserRepository) List(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 查询总数
	if err := r.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.DB.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// FindByCondition 根据条件查询用户列表
func (r *UserRepository) FindByCondition(condition map[string]any, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := r.DB.Model(&models.User{})

	// 应用查询条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// Transaction 执行事务
func (r *UserRepository) Transaction(fn func(*gorm.DB) error) error {
	return r.DB.Transaction(fn)
}
