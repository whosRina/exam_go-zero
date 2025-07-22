package svc

import (
	"exam-system/questionBank/internal/config"
	"exam-system/questionBank/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	QuestionModel     model.QuestionModel
	QuestionBankModel model.QuestionBankModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 这里初始化MySQL连接
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:            c,
		QuestionModel:     model.NewQuestionModel(conn),
		QuestionBankModel: model.NewQuestionBankModel(conn),
	}
}
