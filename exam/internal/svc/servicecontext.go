package svc

import (
	"exam-system/exam/internal/config"
	"exam-system/exam/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	PaperModel          model.PaperModel
	PaperRuleModel      model.PaperRuleModel
	QuestionModel       model.QuestionModel
	ClassModel          model.ClassModel
	ClassMemberModel    model.ClassMemberModel
	ExamModel           model.ExamModel
	ExamAnswerModel     model.ExamAnswerModel
	ExamAttemptModel    model.ExamAttemptModel
	UsersModel          model.UsersModel
	GeneratedPaperModel model.GeneratedPaperModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 这里初始化MySQL连接
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:              c,
		PaperRuleModel:      model.NewPaperRuleModel(conn),
		PaperModel:          model.NewPaperModel(conn),
		QuestionModel:       model.NewQuestionModel(conn),
		ClassModel:          model.NewClassModel(conn),
		ClassMemberModel:    model.NewClassMemberModel(conn),
		ExamModel:           model.NewExamModel(conn),
		ExamAnswerModel:     model.NewExamAnswerModel(conn),
		ExamAttemptModel:    model.NewExamAttemptModel(conn),
		UsersModel:          model.NewUsersModel(conn),
		GeneratedPaperModel: model.NewGeneratedPaperModel(conn),
	}
}
