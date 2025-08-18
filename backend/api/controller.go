package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JsonBack 用于统一返回JSON格式的响应
// 参数说明：
//
//	c       gin的上下文
//	message 返回的消息内容
//	ret     返回码，0表示成功，-1表示服务器错误，-2表示请求参数错误
//	data    返回的数据内容（可选）
func JsonBack(c *gin.Context, message string, ret int, data interface{}) {
	switch ret {
	case 0:
		// 成功返回
		if data != nil {
			// 有数据时返回data字段
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": message,
				"data":    data,
			})
		} else {
			// 无数据时不返回data字段
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": message,
			})
		}
	case -2:
		// 参数错误，返回400
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": message,
		})
	case -1:
		// 服务器错误，返回500
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": message,
		})
	}
}
