package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ExamAnswerModel = (*customExamAnswerModel)(nil)

type (
	// ExamAnswerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customExamAnswerModel.
	ExamAnswerModel interface {
		examAnswerModel
		withSession(session sqlx.Session) ExamAnswerModel
	}

	customExamAnswerModel struct {
		*defaultExamAnswerModel
	}
)

// NewExamAnswerModel returns a model for the database table.
func NewExamAnswerModel(conn sqlx.SqlConn) ExamAnswerModel {
	return &customExamAnswerModel{
		defaultExamAnswerModel: newExamAnswerModel(conn),
	}
}

func (m *customExamAnswerModel) withSession(session sqlx.Session) ExamAnswerModel {
	return NewExamAnswerModel(sqlx.NewSqlConnFromSession(session))
}
