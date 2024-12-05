package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ErrResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Success(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	ctx.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: data,
	})
}

func Err(ctx *gin.Context, code int, msgOrErr interface{}) {
	var msg string
	switch msgOrErr.(type) {
	case string:
		msg = msgOrErr.(string)
	case error:
		err := msgOrErr.(error)
		for {
			if uErr := errors.Unwrap(err); uErr != nil {
				err = uErr
			} else {
				break
			}
		}
		msg = err.Error()
	}
	ctx.AbortWithStatusJSON(code, &ErrResponse{
		Code: code,
		Msg:  msg,
	})
}

func ErrWithData(ctx *gin.Context, code int, data interface{}) {
	ctx.AbortWithStatusJSON(code, &Response{
		Code: code,
		Msg:  "request failed",
		Data: data,
	})
}

func ResourceFromResponseBody(body []byte) (resource interface{}, err error) {
	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("body解析错误:%w", err)
	}
	value, exists := resp["resource"]
	if !exists {
		return nil, fmt.Errorf("resource解析错误:%w", err)
	}
	return value, nil
}
