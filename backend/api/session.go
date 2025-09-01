package api

import(
	"mychat-backend/internal/dto/request"
	"mychat-backend/pkg/zaplog"
	"mychat-backend/pkg/constants"
	"mychat-backend/internal/service/gorm"
	"net/http"
	"github.com/gin-gonic/gin"
)

func OpenSession(c *gin.Context) {
	var session request.OpenSessionRequest
	if err := c.ShouldBind(&session); err != nil {
		zaplog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 500, 
			"message": constants.SYSTEM_ERROR
		})
		return
	}
	msg, 
}