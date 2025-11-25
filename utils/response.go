package utils

import (
	"net/http"

	"github.com/demo/common"
	"github.com/gin-gonic/gin"
)

type CustomResponse interface {
	GetCode() int
	GetMsg() string
	GetData() any
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func (r *Response) GetCode() int {
	return r.Code
}

func (r *Response) GetMsg() string {
	return r.Msg
}

func (r *Response) GetData() any {
	return r.Data
}

// Success 成功返回
func Success(ctx *gin.Context, data any, msg ...string) {
	message := common.GetErrorMessage(common.CodeSuccess)
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}

	ToResponse(ctx, &Response{
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

	ToResponse(ctx, &Response{
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

	ToResponse(ctx, &Response{
		Code: code,
		Msg:  message,
		Data: data,
	})
}

func Error(ctx *gin.Context, err error) {
	if customErr, ok := err.(*common.CustomError); ok {
		ToResponse(ctx, customErr)
	} else {
		ToResponse(ctx, &Response{
			Code: common.CodeInternalError,
			Msg:  "内部错误",
			Data: nil,
		})
	}
}

// SuccessWithMsg 成功返回（自定义消息）
func SuccessWithMsg(ctx *gin.Context, msg string) {
	ToResponse(ctx, &Response{
		Code: common.CodeSuccess,
		Msg:  msg,
	})
}

func ToResponse(ctx *gin.Context, vo CustomResponse) {
	ctx.JSON(http.StatusOK, &Response{
		Code: vo.GetCode(),
		Msg:  vo.GetMsg(),
		Data: vo.GetData(),
	})
}
