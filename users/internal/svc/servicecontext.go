package svc

import (
	"exam-system/users/internal/config"
	"exam-system/users/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	UsersModel model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 这里初始化MySQL连接
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:     c,
		UsersModel: model.NewUsersModel(conn),
	}
}
