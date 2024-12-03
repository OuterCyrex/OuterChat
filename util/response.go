package util

import (
	"github.com/gin-gonic/gin"
)

// swagger:model
type Response struct {
	data interface{}
	code int
	msg  string
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		data: data,
		code: 200,
		msg:  "success",
	}
}

func NewErrorResponse(code int, msg string) *Response {
	return &Response{
		data: nil,
		code: code,
		msg:  msg,
	}
}

func SuccessHttpResponse(data interface{}) *gin.H {
	r := NewSuccessResponse(data)
	return &gin.H{
		"code": r.code,
		"msg":  r.msg,
		"data": r.data,
	}
}

func ErrorHttpResponse(code int, msg string) *gin.H {
	r := NewErrorResponse(code, msg)
	return &gin.H{
		"code": r.code,
		"msg":  r.msg,
		"data": r.data,
	}
}

func SendErrorResponse(c *gin.Context, code int, msg string) {
	c.JSON(-1, ErrorHttpResponse(code, msg))
	c.Abort()
}

func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(-1, SuccessHttpResponse(data))
	c.Abort()
}
