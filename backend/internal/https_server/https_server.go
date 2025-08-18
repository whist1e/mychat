package https_server

import (
	api "mychat-backend/api"
	"mychat-backend/internal/config"

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

	// 注意：SSL/TLS处理已移除，当前为HTTP服务器
	// 如需HTTPS，请配置SSL证书和TLS处理中间件
	// GE.Use(ssl.TlsHandler(config.GetConfig().MainConfig.Host, config.GetConfig().MainConfig.Port))

	// 静态文件服务
	GE.Static("/static/avatars", config.GetConfig().StaticSrcConfig.StaticAvatarPath)
	GE.Static("/static/files", config.GetConfig().StaticSrcConfig.StaticFilePath)

	// API路由
	GE.POST("/login", api.Login)
	GE.POST("/register", api.Register)

	// 其他API路由（已注释）
	// GE.POST("/user/updateUserInfo", api.UpdateUserInfo)
	// GE.POST("/user/getUserInfoList", api.GetUserInfoList)
	// GE.POST("/user/ableUsers", api.AbleUsers)
	// GE.POST("/user/getUserInfo", api.GetUserInfo)
	// GE.POST("/user/disableUsers", api.DisableUsers)
	// GE.POST("/user/deleteUsers", api.DeleteUsers)
	// GE.POST("/user/setAdmin", api.SetAdmin)
	// GE.POST("/user/sendSmsCode", api.SendSmsCode)
	// GE.POST("/user/smsLogin", api.SmsLogin)
	// GE.POST("/user/wsLogout", api.WsLogout)
	// GE.POST("/group/createGroup", api.CreateGroup)
	// GE.POST("/group/loadMyGroup", api.LoadMyGroup)
	// GE.POST("/group/checkGroupAddMode", api.CheckGroupAddMode)
	// GE.POST("/group/enterGroupDirectly", api.EnterGroupDirectly)
	// GE.POST("/group/leaveGroup", api.LeaveGroup)
	// GE.POST("/group/dismissGroup", api.DismissGroup)
	// GE.POST("/group/getGroupInfo", api.GetGroupInfo)
	// GE.POST("/group/getGroupInfoList", api.GetGroupInfoList)
	// GE.POST("/group/deleteGroups", api.DeleteGroups)
	// GE.POST("/group/setGroupsStatus", api.SetGroupsStatus)
	// GE.POST("/group/updateGroupInfo", api.UpdateGroupInfo)
	// GE.POST("/group/getGroupMemberList", api.GetGroupMemberList)
	// GE.POST("/group/removeGroupMembers", api.RemoveGroupMembers)
	// GE.POST("/session/openSession", api.OpenSession)
	// GE.POST("/session/getUserSessionList", api.GetUserSessionList)
	// GE.POST("/session/getGroupSessionList", api.GetGroupSessionList)
	// GE.POST("/session/deleteSession", api.DeleteSession)
	// GE.POST("/session/checkOpenSessionAllowed", api.CheckOpenSessionAllowed)
	// GE.POST("/contact/getUserList", api.GetUserList)
	// GE.POST("/contact/loadMyJoinedGroup", api.LoadMyJoinedGroup)
	// GE.POST("/contact/getContactInfo", api.GetContactInfo)
	// GE.POST("/contact/deleteContact", api.DeleteContact)
	// GE.POST("/contact/applyContact", api.ApplyContact)
	// GE.POST("/contact/getNewContactList", api.GetNewContactList)
	// GE.POST("/contact/passContactApply", api.PassContactApply)
	// GE.POST("/contact/blackContact", api.BlackContact)
	// GE.POST("/contact/cancelBlackContact", api.CancelBlackContact)
	// GE.POST("/contact/getAddGroupList", api.GetAddGroupList)
	// GE.POST("/contact/refuseContactApply", api.RefuseContactApply)
	// GE.POST("/contact/blackApply", api.BlackApply)
	// GE.POST("/message/getMessageList", api.GetMessageList)
	// GE.POST("/message/getGroupMessageList", api.GetGroupMessageList)
	// GE.POST("/message/uploadAvatar", api.UploadAvatar)
	// GE.POST("/message/uploadFile", api.UploadFile)
	// GE.POST("/chatroom/getCurContactListInChatRoom", api.GetCurContactListInChatRoom)
	// GE.GET("/wss", api.WsLogin)
}
