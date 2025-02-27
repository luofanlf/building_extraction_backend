package dto

import "github.com/gin-gonic/gin"

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"err"`
}

// OkWithData 返回成功响应和数据
func OkWithData(data interface{}, c *gin.Context) {
	c.JSON(200, Response{
		Code:  0,
		Data:  data,
		Error: "success",
	})
}

// FailWithMessage 返回失败响应和错误消息
func FailWithMessage(message string, c *gin.Context) {
	c.JSON(400, Response{
		Code:  -1,
		Error: message,
	})
}
