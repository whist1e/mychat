package api

import (
	"mychat-backend/internal/dto/request"
	"mychat-backend/internal/dto/respond"
	"mychat-backend/pkg/zaplog"
	"mychat-backend/service/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, respond.BadRequestResponse("参数绑定失败"))
		return
	}

	service := gorm.NewUserInfoService()
	message, userInfo, err := service.Register(req)
	if err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, respond.ErrorResponse(500, "注册失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, respond.SuccessResponse(userInfo, message))
}

func Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, respond.BadRequestResponse("参数绑定失败"))
		return
	}

	service := gorm.NewUserInfoService()
	message, userInfo, err := service.Login(req)
	if err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, respond.ErrorResponse(500, "登录失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, respond.SuccessResponse(userInfo, message))
}
