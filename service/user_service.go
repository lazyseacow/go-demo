package service

import (
	"github.com/demo/common"
	"github.com/demo/models"
	"github.com/demo/repository"
	"github.com/demo/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	*BaseService
	repo repository.UserRepositoryInterface
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		BaseService: NewBaseService(),
		repo:        repository.NewUserRepository(),
	}
}

// NewUserServiceWithRepo 创建用户服务实例（用于测试）
func NewUserServiceWithRepo(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		BaseService: NewBaseService(),
		repo:        repo,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

// Register 用户注册
func (s *UserService) Register(req RegisterRequest) (*models.User, error) {
	// 检查用户名是否存在
	exists, err := s.repo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, common.NewError(common.CodeDBQueryFailed, "查询用户失败")
	}
	if exists {
		return nil, common.NewError(common.CodeUsernameExists)
	}

	// 检查邮箱是否存在
	exists, err = s.repo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, common.NewError(common.CodeDBQueryFailed, "查询邮箱失败")
	}
	if exists {
		return nil, common.NewError(common.CodeEmailExists)
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.NewError(common.CodeInternalError, "密码加密失败")
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Status:   1,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, common.NewError(common.CodeDBInsertFailed, "创建用户失败")
	}

	return user, nil
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string       `json:"token"`
	UserInfo *models.User `json:"user_info"`
}

// Login 用户登录
func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
	// 查询用户
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewError(common.CodeInvalidPassword)
		}
		return nil, common.NewError(common.CodeDBQueryFailed, "查询用户失败")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, common.NewError(common.CodeInvalidPassword)
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, common.NewError(common.CodeUserDisabled)
	}

	// 生成 Token
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return nil, common.NewError(common.CodeTokenGenFailed)
	}

	return &LoginResponse{
		Token:    token,
		UserInfo: user,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID int64) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewError(common.CodeUserNotFound)
		}
		return nil, common.NewError(common.CodeDBQueryFailed, "查询用户失败")
	}

	return user, nil
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(page, pageSize int) (*models.PageResponse, error) {
	// 默认分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, common.NewError(common.CodeDBQueryFailed, "查询用户列表失败")
	}

	return &models.PageResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     users,
	}, nil
}

// GetUserByID 根据 ID 获取用户
func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewError(common.CodeUserNotFound)
		}
		return nil, common.NewError(common.CodeDBQueryFailed, "查询用户失败")
	}

	return user, nil
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Phone  string `json:"phone"`
	Avatar string `json:"avatar"`
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID int64, req UpdateUserRequest) error {
	// 构建更新字段
	updates := make(map[string]any)
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) == 0 {
		return common.NewError(common.CodeParamInvalid, "没有需要更新的字段")
	}

	if err := s.repo.UpdateFields(userID, updates); err != nil {
		return common.NewError(common.CodeDBUpdateFailed, "更新用户失败")
	}

	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id, currentUserID int64) error {
	// 不能删除自己
	if id == currentUserID {
		return common.NewError(common.CodeCannotDeleteSelf)
	}

	// 检查用户是否存在
	_, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.NewError(common.CodeUserNotFound)
		}
		return common.NewError(common.CodeDBQueryFailed, "查询用户失败")
	}

	// 删除用户
	if err := s.repo.Delete(id); err != nil {
		return common.NewError(common.CodeDBDeleteFailed, "删除用户失败")
	}

	return nil
}
