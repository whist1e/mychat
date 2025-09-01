package gorm

import (
	"encoding/json"
	"errors"
	"mychat-backend/internal/dao"
	"mychat-backend/internal/dto/respond"
	"mychat-backend/internal/model"
	myredis "mychat-backend/internal/service/redis"
	"mychat-backend/pkg/constants"
	"mychat-backend/pkg/enum/contact/contact_type_enum"
	"mychat-backend/pkg/zaplog"
	"time"

	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"
)

type userContactService struct{}

var UserContactService = new(userContactService)

// 获取联系人列表
func (userContactService) GetUserList(ownerId string) (string, []respond.MyUserListRespond, int) {
	rspString, err := myredis.GetKeyNilIsErr("contact_user_list" + ownerId)
	if err != nil {
		//缓存未命中
		if errors.Is(err, redis.Nil) {
			//准备在数据中查找contact，但是其中包含了群组
			var contactList []model.UserContact
			//不存在不属于错误，所以用info
			//查找未被删除的联系人
			//使用软删除，使用status记录状态
			if res := dao.GormDB.Order("created_at DESC").Where("user_id = ? AND status != 4", ownerId).Find(&contactList); res.Error != nil {
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					msg := "目前不存在此联系人"
					zaplog.Info(msg)
					return msg, nil, 0
				} else {
					zaplog.Error(res.Error.Error())
					return constants.SYSTEM_ERROR, nil, -1
				}
			}

			//从contact中查找联系人
			var userList []respond.MyUserListRespond
			for _, contact := range contactList {
				if contact.ContactType == contact_type_enum.USER {
					var user model.UserInfo
					if res := dao.GormDB.First(&user, "uuid = ?", contact.ContactId); res.Error != nil {
						zaplog.Error(res.Error.Error())
						return constants.SYSTEM_ERROR, nil, -1
					}
					userList = append(userList, respond.MyUserListRespond{
						UserId:   user.Uuid,
						UserName: user.Nickname,
						Avatar:   user.Avatar,
					})
				}
			}

			//放到缓存
			rspString, err := json.Marshal(userList)
			if err != nil {
				zaplog.Error(err.Error())
			}
			if myredis.SetKeyEx("contact_user_list"+ownerId, string(rspString), time.Minute*constants.REdis_TIMEOUT); err != nil {
				zaplog.Error(err.Error())
			}
			return "获取列表成功", userList, 0
		} else {
			zaplog.Error(err.Error())
		}
	}
	var rsp []respond.MyUserListRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zaplog.Error(err.Error())
	}
	//test缓存服务
	// zaplog.Info("已成功缓存")
	return "获取列表成功", rsp, 0
}
