package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误码定义
const (
	// CodeSuccess ==================== 通用错误码 (1000-1999) ====================
	CodeSuccess          = 200 // 成功
	CodeBadRequest       = 400 // 请求参数错误
	CodeUnauthorized     = 401 // 未授权
	CodeForbidden        = 403 // 禁止访问
	CodeNotFound         = 404 // 资源不存在
	CodeMethodNotAllowed = 405 // 方法不允许
	CodeConflict         = 409 // 资源冲突
	CodeTooManyRequests  = 429 // 请求过于频繁
	CodeInternalError    = 500 // 服务器内部错误
	CodeServiceUnavail   = 503 // 服务不可用

	// CodeTokenMissing ==================== 认证相关错误码 (10000-10999) ====================
	CodeTokenMissing   = 10001 // Token 缺失
	CodeTokenInvalid   = 10002 // Token 无效
	CodeTokenExpired   = 10003 // Token 已过期
	CodeTokenGenFailed = 10004 // Token 生成失败
	CodeLoginRequired  = 10005 // 需要登录

	// CodeUserExists ==================== 用户相关错误码 (11000-11999) ====================
	CodeUserExists       = 11001 // 用户已存在
	CodeUserNotFound     = 11002 // 用户不存在
	CodeUsernameExists   = 11003 // 用户名已存在
	CodeEmailExists      = 11004 // 邮箱已被注册
	CodePhoneExists      = 11005 // 手机号已被注册
	CodeInvalidPassword  = 11006 // 密码错误
	CodePasswordTooWeak  = 11007 // 密码强度不够
	CodeUserDisabled     = 11008 // 账号已被禁用
	CodeUserDeleted      = 11009 // 账号已被删除
	CodeCannotDeleteSelf = 11010 // 不能删除自己

	// CodeArticleNotFound ==================== 文章相关错误码 (12000-12999) ====================
	CodeArticleNotFound   = 12001 // 文章不存在
	CodeArticleExists     = 12002 // 文章已存在
	CodeArticleCreateFail = 12003 // 文章创建失败
	CodeArticleUpdateFail = 12004 // 文章更新失败
	CodeArticleDeleteFail = 12005 // 文章删除失败
	CodeNoPermission      = 12006 // 没有权限操作
	CodeInvalidArticleID  = 12007 // 无效的文章ID

	// CodeDBError ==================== 数据库相关错误码 (13000-13999) ====================
	CodeDBError         = 13001 // 数据库错误
	CodeDBConnectFailed = 13002 // 数据库连接失败
	CodeDBQueryFailed   = 13003 // 数据库查询失败
	CodeDBInsertFailed  = 13004 // 数据库插入失败
	CodeDBUpdateFailed  = 13005 // 数据库更新失败
	CodeDBDeleteFailed  = 13006 // 数据库删除失败
	CodeRedisError      = 13007 // Redis 错误
	CodeMongoDBError    = 13008 // MongoDB 错误

	// CodeParamInvalid ==================== 参数验证错误码 (14000-14999) ====================
	CodeParamInvalid      = 14001 // 参数无效
	CodeParamMissing      = 14002 // 参数缺失
	CodeParamTypeMismatch = 14003 // 参数类型错误
	CodeParamOutOfRange   = 14004 // 参数超出范围
	CodeInvalidEmail      = 14005 // 邮箱格式错误
	CodeInvalidPhone      = 14006 // 手机号格式错误
	CodeInvalidUsername   = 14007 // 用户名格式错误

	// CodeOperationFailed ==================== 业务逻辑错误码 (15000-15999) ====================
	CodeOperationFailed = 15001 // 操作失败
	CodeDataExists      = 15002 // 数据已存在
	CodeDataNotFound    = 15003 // 数据不存在
	CodeStatusError     = 15004 // 状态错误
	CodeLimitExceeded   = 15005 // 超出限制
)

// 错误码对应的错误消息
var errMessages = map[int]string{
	// 通用错误
	CodeSuccess:          "操作成功",
	CodeBadRequest:       "请求参数错误",
	CodeUnauthorized:     "未授权，请先登录",
	CodeForbidden:        "禁止访问",
	CodeNotFound:         "资源不存在",
	CodeMethodNotAllowed: "方法不允许",
	CodeConflict:         "资源冲突",
	CodeTooManyRequests:  "请求过于频繁，请稍后再试",
	CodeInternalError:    "服务器内部错误",
	CodeServiceUnavail:   "服务暂时不可用",

	// 认证相关
	CodeTokenMissing:   "Token 缺失",
	CodeTokenInvalid:   "Token 无效",
	CodeTokenExpired:   "Token 已过期，请重新登录",
	CodeTokenGenFailed: "Token 生成失败",
	CodeLoginRequired:  "请先登录",

	// 用户相关
	CodeUserExists:       "用户已存在",
	CodeUserNotFound:     "用户不存在",
	CodeUsernameExists:   "用户名已存在",
	CodeEmailExists:      "邮箱已被注册",
	CodePhoneExists:      "手机号已被注册",
	CodeInvalidPassword:  "用户名或密码错误",
	CodePasswordTooWeak:  "密码强度不够",
	CodeUserDisabled:     "账号已被禁用，请联系管理员",
	CodeUserDeleted:      "账号已被删除",
	CodeCannotDeleteSelf: "不能删除自己",

	// 文章相关
	CodeArticleNotFound:   "文章不存在",
	CodeArticleExists:     "文章已存在",
	CodeArticleCreateFail: "文章创建失败",
	CodeArticleUpdateFail: "文章更新失败",
	CodeArticleDeleteFail: "文章删除失败",
	CodeNoPermission:      "没有权限操作该资源",
	CodeInvalidArticleID:  "无效的文章ID",

	// 数据库相关
	CodeDBError:         "数据库错误",
	CodeDBConnectFailed: "数据库连接失败",
	CodeDBQueryFailed:   "数据库查询失败",
	CodeDBInsertFailed:  "数据库插入失败",
	CodeDBUpdateFailed:  "数据库更新失败",
	CodeDBDeleteFailed:  "数据库删除失败",
	CodeRedisError:      "缓存服务错误",
	CodeMongoDBError:    "文档数据库错误",

	// 参数验证
	CodeParamInvalid:      "参数无效",
	CodeParamMissing:      "缺少必要参数",
	CodeParamTypeMismatch: "参数类型错误",
	CodeParamOutOfRange:   "参数超出允许范围",
	CodeInvalidEmail:      "邮箱格式错误",
	CodeInvalidPhone:      "手机号格式错误",
	CodeInvalidUsername:   "用户名格式错误",

	// 业务逻辑
	CodeOperationFailed: "操作失败",
	CodeDataExists:      "数据已存在",
	CodeDataNotFound:    "数据不存在",
	CodeStatusError:     "状态错误",
	CodeLimitExceeded:   "超出限制",
}

// GetErrorMessage 根据错误码获取错误消息
func GetErrorMessage(code int) string {
	if msg, ok := errMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data,omitempty"`
}

// Success 成功返回
func Success(ctx *gin.Context, data any, msg ...string) {
	message := GetErrorMessage(CodeSuccess)
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}

	ctx.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  message,
		Data: data,
	})
}

// Fail 失败返回（使用错误码）
func Fail(ctx *gin.Context, code int, msg ...string) {
	message := GetErrorMessage(code)
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}

	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
	})
}

// FailWithData 失败返回（带数据）
func FailWithData(ctx *gin.Context, code int, data any, msg ...string) {
	message := GetErrorMessage(code)
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}

	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
		Data: data,
	})
}

// Error 返回自定义错误
func Error(ctx *gin.Context, err *Response) {
	ctx.JSON(http.StatusOK, Response{
		Code: err.Code,
		Msg:  err.Msg,
		Data: err.Data,
	})
}

//func (e *Response) Error() string {
//	return fmt.Sprintf("code: %d, msg: %s, data: %+v", e.Code, e.Msg, e.Data)
//}

func NewError(code int, message ...string) *Response {
	msg := GetErrorMessage(code)
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return &Response{
		Code: code,
		Msg:  msg,
	}
}
