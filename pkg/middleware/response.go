package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseCode 错误码类型，1000以下为通用码，1000以上为用户自定义码
type ResponseCode int

const (
	SuccessCode       ResponseCode = iota // 成功
	UndefErrorCode                        // 未定义的错误
	ValidErrorCode                        // 校验错误
	InternalErrorCode                     // 内部错误

	InvalidRequestErrorCode ResponseCode = 401  // 请求未经授权
	CustomizeCode           ResponseCode = 1000 // 自定义错误码

	GROUPALL_SAVE_FLOWERROR ResponseCode = 2001 // 自定义错误码
)

// Response 响应结构体
type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      interface{}  `json:"data"`
	TraceID   interface{}  `json:"trace_id"`
	Stack     interface{}  `json:"stack"`
}

// ResponseError 错误响应函数，code 为错误码，err 为错误信息
func ResponseError(c *gin.Context, code ResponseCode, err error) {
	// 设置响应头
	c.Header("Content-Type", "application/json; charset=utf-8")
	// 构造响应体
	c.JSON(http.StatusOK, Response{
		ErrorCode: code,
		ErrorMsg:  err.Error(),
		Data:      nil,
		TraceID:   c.GetHeader("Trace-Id"),
		Stack:     err.Error(),
	})
}

// ResponseSuccess 成功响应函数，data 为响应数据
func ResponseSuccess(c *gin.Context, data interface{}) {
	// 设置响应头
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 构造响应体
	c.JSON(http.StatusOK, Response{
		ErrorCode: SuccessCode,
		ErrorMsg:  "成功",
		Data:      data,
		TraceID:   c.GetHeader("Trace-Id"),
		Stack:     nil,
	})
}
