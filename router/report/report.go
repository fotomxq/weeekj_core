package RouterReport

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//常见的头数据结构封装

// 通用数据反馈头
type DataType struct {
	//错误信息
	Status bool `json:"status"`
	//错误信息
	Code string `json:"code"`
	//错误描述
	Msg string `json:"msg"`
	//数据个数
	Count int64 `json:"count"`
	//数据集合
	Data interface{} `json:"data"`
}

// 反馈成功
func BaseSuccess(c *gin.Context) {
	res := DataType{
		Status: true,
		Code:   "",
		Msg:    "",
		Count:  0,
		Data:   nil,
	}
	c.JSON(http.StatusOK, &res)
	c.Abort()
}

// 反馈成功或失败
func BaseBool(c *gin.Context, code string, b bool, errMsg string) {
	if b {
		BaseSuccess(c)
	} else {
		BaseError(c, code, errMsg)
	}
}

// 反馈一般数据
func BaseData(c *gin.Context, data interface{}) {
	//反馈数据
	res := DataType{
		Status: true,
		Code:   "",
		Msg:    "",
		Count:  0,
		Data:   data,
	}
	c.JSON(http.StatusOK, &res)
	c.Abort()
}

// 反馈列队数据
func BaseDataList(c *gin.Context, count int64, data interface{}) {
	//反馈数据
	res := DataType{
		Status: true,
		Code:   "",
		Msg:    "",
		Count:  count,
		Data:   data,
	}
	c.JSON(http.StatusOK, &res)
	c.Abort()
}

// 反馈错误
func BaseError(c *gin.Context, code string, msg string) {
	res := DataType{
		false,
		code,
		msg,
		0,
		nil,
	}
	c.JSON(http.StatusOK, &res)
	c.Abort()
}

// BaseCustomMsg 反馈自定义内容
func BaseCustomMsg(c *gin.Context, status bool, code string, msg string) {
	res := DataType{
		status,
		code,
		msg,
		0,
		nil,
	}
	c.JSON(http.StatusOK, &res)
	c.Abort()
}
