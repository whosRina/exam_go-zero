package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ExamModel = (*customExamModel)(nil)

type (
	// ExamModel is an interface to be customized, add more methods here,
	// and implement the added methods in customExamModel.
	ExamModel interface {
		examModel
		withSession(session sqlx.Session) ExamModel
	}

	customExamModel struct {
		*defaultExamModel
	}
)

// NewExamModel returns a model for the database table.
func NewExamModel(conn sqlx.SqlConn) ExamModel {
	return &customExamModel{
		defaultExamModel: newExamModel(conn),
	}
}

func (m *customExamModel) withSession(session sqlx.Session) ExamModel {
	return NewExamModel(sqlx.NewSqlConnFromSession(session))
}
