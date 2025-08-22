package api

import (
	"mychat-backend/internal/dto/request"
	"mychat-backend/pkg/zaplog"
	"mychat-backend/pkg/constants"
	"mychat-backend/service/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}

	message, userInfo, ret := gorm.UserInfoService.Register(req)
	JsonBack(c, message, ret, userInfo)	
}

//当前没有加短信验证功能，只能通过密码登录
func Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}

	message, userInfo, ret := gorm.UserInfoService.PasswordLogin(req.Telephone, req.Password)
	JsonBack(c, message, ret, userInfo)
}

func UpdateUserInfo(c *gin.Context) {
	var req request.GetUserInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
}

func GetUserInfoList(c *gin.Context) {
	var req request.GetUserInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, userInfo, ret := gorm.UserInfoService.GetUserInfo(req.Uuid)
	JsonBack(c, message, ret, userInfo)
}
