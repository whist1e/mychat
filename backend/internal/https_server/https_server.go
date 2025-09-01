package https_server

import (
	api "mychat-backend/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var GE *gin.Engine

func init() {
	GE = gin.Default()

	// 配置CORS跨域
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	GE.Use(cors.New(corsConfig))

	// API路由
	GE.POST("/user/login", api.Login)
	GE.POST("/user/register", api.Register)
	GE.POST("/user/getUserInfo", api.GetUserInfo)
	GE.POST("/user/updateUserInfo", api.UpdateUserInfo)

	GE.POST("/contact/getUserList", api.GetUserList)

	GE.POST("/session/openSession", api.OpenSession)
}
