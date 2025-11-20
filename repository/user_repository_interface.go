package repository

import "github.com/demo/models"

// UserRepositoryInterface 用户仓库接口
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByID(id int64) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Update(user *models.User) error
	UpdateFields(id int64, fields map[string]any) error
	Delete(id int64) error
	List(page, pageSize int) ([]models.User, int64, error)
}
