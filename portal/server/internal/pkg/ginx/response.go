package ginx

import (
	"errors"
	"net/http"
	"portal/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

const (
	ResponseSuccess = 0
	ResponseFailure = 500
)

type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

// SuccessResponse 正常响应
func SuccessResponse(context *gin.Context, data interface{}) {
	response := Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}

	// 添加请求ID
	if requestID, exists := context.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	context.JSON(http.StatusOK, response)
}

// FailureResponse 异常响应
func FailureResponse(context *gin.Context, err error) {
	if err == nil {
		return
	}

	var code = errorx.UnknownErrorCode
	var msg = err.Error()
	var bizErr *errorx.Error
	if errors.As(err, &bizErr) {
		code = bizErr.Code
		msg = bizErr.Message
	}

	response := Response{
		Code: code,
		Msg:  msg,
	}

	// 添加请求ID
	if requestID, exists := context.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	context.AbortWithStatusJSON(http.StatusOK, response)
}
func FailResponseWithCode(context *gin.Context, code int, msg string) {
	response := Response{
		Code: code,
		Msg:  msg,
	}

	// 添加请求ID
	if requestID, exists := context.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	context.AbortWithStatusJSON(http.StatusOK, response)
}
