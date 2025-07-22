package svc

import (
	"exam-system/classes/internal/config"
	"exam-system/classes/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	ClassModel       model.ClassModel
	ClassMemberModel model.ClassMemberModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 这里初始化MySQL连接
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:           c,
		ClassModel:       model.NewClassModel(conn),
		ClassMemberModel: model.NewClassMemberModel(conn),
	}
}
