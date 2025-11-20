package controllers

import (
	"github.com/demo/common"
	"github.com/demo/service"
	"github.com/demo/utils"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	*BaseController
	userService *service.UserService
}

// NewAuthController 创建认证控制器实例
func NewAuthController() *AuthController {
	return &AuthController{
		BaseController: NewBaseController(),
		userService:    service.NewUserService(),
	}
}

// Register 用户注册
// @Summary      用户注册
// @Description  注册新用户账号
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      service.RegisterRequest  true  "注册参数"
// @Success      200      {object}  utils.Response{data=object{user_id=int,username=string}}  "注册成功"
// @Failure      14001    {object}  utils.Response  "参数错误"
// @Failure      11003    {object}  utils.Response  "用户名已存在"
// @Failure      11004    {object}  utils.Response  "邮箱已被注册"
// @Router       /auth/register [post]
func (ctrl *AuthController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Fail(ctx, common.CodeParamInvalid, "参数错误: "+err.Error())
		return
	}

	// 调用 Service 层处理业务逻辑
	user, err := ctrl.userService.Register(req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"user_id":  user.ID,
		"username": user.Username,
	}, "注册成功")
}

// Login 用户登录
// @Summary      用户登录
// @Description  使用用户名和密码登录系统
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      service.LoginRequest  true  "登录参数"
// @Success      200      {object}  utils.Response{data=service.LoginResponse}  "登录成功"
// @Failure      14001    {object}  utils.Response  "参数错误"
// @Failure      11006    {object}  utils.Response  "用户名或密码错误"
// @Failure      11008    {object}  utils.Response  "账号已被禁用"
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Fail(ctx, common.CodeParamInvalid, "参数错误: "+err.Error())
		return
	}

	// 调用 Service 层处理业务逻辑
	resp, err := ctrl.userService.Login(req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, resp, "登录成功")
}

// Logout 用户登出
// @Summary      用户登出
// @Description  退出登录，清除登录状态
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  utils.Response  "登出成功"
// @Router       /auth/logout [post]
func (ctrl *AuthController) Logout(ctx *gin.Context) {
	// TODO: 可以在 Redis 中维护一个 Token 黑名单
	utils.Success(ctx, nil, "登出成功")
}
