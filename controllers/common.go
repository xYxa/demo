package controllers

import "github.com/gin-gonic/gin"

type JsonStruct struct {
	Code  int
	Msg   interface{}
	Data  interface{}
	Count int64
}
type JsonErrStruct struct {
	Code int
	Msg  interface{}
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, data interface{}, count int64) {
	json := &JsonStruct{Code: code, Msg: msg, Data: data, Count: count}
	c.JSON(200, json)
}

func ReturnError(c *gin.Context, code int, msg interface{}) {
	json := &JsonErrStruct{Code: code, Msg: msg}
	c.JSON(200, json)
}
