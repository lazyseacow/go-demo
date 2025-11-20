package utils

import (
	"net/http"

	"github.com/demo/common"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Success 成功返回
func Success(ctx *gin.Context, data any, msg ...string) {
	message := common.GetErrorMessage(common.CodeSuccess)
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}

	ctx.JSON(http.StatusOK, Response{
		Code: common.CodeSuccess,
		Msg:  message,
		Data: data,
	})
}

// Fail 失败返回（使用错误码）
func Fail(ctx *gin.Context, code int, msg ...string) {
	message := common.GetErrorMessage(code)
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
	message := common.GetErrorMessage(code)
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
func Error(ctx *gin.Context, err *common.CustomError) {
	ctx.JSON(http.StatusOK, Response{
		Code: err.Code,
		Msg:  err.Message,
		Data: err.Data,
	})
}

// SuccessWithMsg 成功返回（自定义消息）
func SuccessWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, Response{
		Code: common.CodeSuccess,
		Msg:  msg,
	})
}
