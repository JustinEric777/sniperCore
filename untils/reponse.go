package untils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type responseJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Ts   string      `json:"ts"`
}

func JsonSuccess(ctx *gin.Context, data interface{}) {
	var response *responseJson

	switch data.(type) {
	case string:
		response = &responseJson{
			Code: 0,
			Msg:  fmt.Sprintf("%v", data),
			Data: gin.H{},
			Ts:   fmt.Sprintf("%d", time.Now().UnixNano()/1e6),
		}
	default:
		response = &responseJson{
			Code: 0,
			Msg:  "",
			Data: data,
			Ts:   fmt.Sprintf("%d", time.Now().UnixNano()/1e6),
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func JsonError(ctx *gin.Context, code int, msg string) {
	response := &responseJson{
		Code: code,
		Msg:  msg,
		Data: "",
		Ts:   fmt.Sprintf("%d", time.Now().UnixNano()/1e6),
	}

	var httpCode int
	if code == 400 {
		httpCode = http.StatusBadRequest
	} else if code == 500 {
		httpCode = http.StatusInternalServerError
	} else {
		httpCode = http.StatusInternalServerError
	}

	ctx.JSON(httpCode, response)
}
