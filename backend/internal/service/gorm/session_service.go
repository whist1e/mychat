package gorm

import (
	"encoding/json"
	"errors"
	"mychat-backend/internal/dao"
	"mychat-backend/internal/dto/request"
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

type sessionService struct{}

var SessionService = new(sessionService)

func (sessionService)OpenSession(r *request.OpenSessionRequest) (string, string, int){
	session := model.Session
	if res := dao.GormDB.Where("send_id = ? and receive_id = ?", r.SendId, r.ReceiveId); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			dao.GormDB.Create()
		}
	}
}