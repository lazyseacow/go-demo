package controllers

import (
	"strconv"

	"github.com/demo/common"
	"github.com/demo/middleware"
	"github.com/demo/models"
	"github.com/demo/service"
	"github.com/demo/utils"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	*BaseController
	userService *service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		BaseController: NewBaseController(),
		userService:    service.NewUserService(),
	}
}

// GetUserList 获取用户列表（分页）
// @Summary      获取用户列表
// @Description  分页获取所有用户列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        page       query     int  false  "页码"  default(1)  minimum(1)
// @Param        page_size  query     int  false  "每页数量"  default(10)  minimum(1)  maximum(100)
// @Success      200        {object}  utils.Response{data=models.PageResponse}  "获取成功"
// @Failure      10005      {object}  utils.Response  "需要登录"
// @Router       /users [get]
func (ctrl *UserController) GetUserList(ctx *gin.Context) {
	var req models.PageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req = models.DefaultPageRequest()
	}

	// 调用 Service 层
	result, err := ctrl.userService.GetUserList(req.Page, req.PageSize)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, result)
}

// GetUserByID 根据 ID 获取用户
// @Summary      获取指定用户
// @Description  根据用户 ID 获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  utils.Response{data=models.User}  "获取成功"
// @Failure      14001  {object}  utils.Response  "无效的用户 ID"
// @Failure      11002  {object}  utils.Response  "用户不存在"
// @Router       /users/{id} [get]
func (ctrl *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Fail(ctx, common.CodeParamInvalid, "无效的用户 ID")
		return
	}

	// 调用 Service 层
	user, err := ctrl.userService.GetUserByID(id)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, user)
}

// UpdateUser 更新用户信息
// @Summary      更新用户信息
// @Description  更新当前登录用户的个人信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      service.UpdateUserRequest  true  "更新参数"
// @Success      200      {object}  utils.Response  "更新成功"
// @Failure      14001    {object}  utils.Response  "参数错误"
// @Failure      10005    {object}  utils.Response  "需要登录"
// @Router       /users/update [post]
func (ctrl *UserController) UpdateUser(ctx *gin.Context) {
	// 获取当前登录用户 ID
	currentUserID := middleware.GetUserID(ctx)
	if currentUserID == 0 {
		utils.Fail(ctx, common.CodeLoginRequired)
		return
	}

	var req service.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Fail(ctx, common.CodeParamInvalid, "参数错误: "+err.Error())
		return
	}

	// 调用 Service 层
	if err := ctrl.userService.UpdateUser(currentUserID, req); err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, nil, "更新成功")
}

// DeleteUser 删除用户（软删除）
// @Summary      删除用户
// @Description  删除指定用户（软删除，可恢复）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "用户ID"
// @Success      200    {object}  utils.Response  "删除成功"
// @Failure      14001  {object}  utils.Response  "无效的用户 ID"
// @Failure      11010  {object}  utils.Response  "不能删除自己"
// @Failure      11002  {object}  utils.Response  "用户不存在"
// @Router       /users/{id}/delete [post]
func (ctrl *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Fail(ctx, common.CodeParamInvalid, "无效的用户 ID")
		return
	}

	// 获取当前登录用户 ID
	currentUserID := middleware.GetUserID(ctx)

	// 调用 Service 层
	if err := ctrl.userService.DeleteUser(id, currentUserID); err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, nil, "删除成功")
}
