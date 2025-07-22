package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ExamQuestionModel = (*customExamQuestionModel)(nil)

type (
	// ExamQuestionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customExamQuestionModel.
	ExamQuestionModel interface {
		examQuestionModel
		withSession(session sqlx.Session) ExamQuestionModel
	}

	customExamQuestionModel struct {
		*defaultExamQuestionModel
	}
)

// NewExamQuestionModel returns a model for the database table.
func NewExamQuestionModel(conn sqlx.SqlConn) ExamQuestionModel {
	return &customExamQuestionModel{
		defaultExamQuestionModel: newExamQuestionModel(conn),
	}
}

func (m *customExamQuestionModel) withSession(session sqlx.Session) ExamQuestionModel {
	return NewExamQuestionModel(sqlx.NewSqlConnFromSession(session))
}
