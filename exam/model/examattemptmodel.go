package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ExamAttemptModel = (*customExamAttemptModel)(nil)

type (
	// ExamAttemptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customExamAttemptModel.
	ExamAttemptModel interface {
		examAttemptModel
		withSession(session sqlx.Session) ExamAttemptModel
	}

	customExamAttemptModel struct {
		*defaultExamAttemptModel
	}
)

// NewExamAttemptModel returns a model for the database table.
func NewExamAttemptModel(conn sqlx.SqlConn) ExamAttemptModel {
	return &customExamAttemptModel{
		defaultExamAttemptModel: newExamAttemptModel(conn),
	}
}

func (m *customExamAttemptModel) withSession(session sqlx.Session) ExamAttemptModel {
	return NewExamAttemptModel(sqlx.NewSqlConnFromSession(session))
}
