package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data,omitempty"`
}

func (r *Response[T]) Success(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

func (r *Response[T]) Fail(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, r)
}

func Resp[T any](code int, msg string, data T) *Response[T] {
	return &Response[T]{Code: code, Msg: msg, Data: data}
}
