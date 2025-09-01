package api

import (
	"mychat-backend/internal/dto/request"
	"mychat-backend/pkg/zaplog"
	"mychat-backend/pkg/constants"
	"mychat-backend/internal/service/gorm"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
	var myUserList request.MyUserListRequest
	if err := c.ShouldBind(&myUserList); err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	msg, userList, ret := gorm.UserContactService.GetUserList(myUserList.UserId)
	JsonBack(c, msg, ret, userList)
}