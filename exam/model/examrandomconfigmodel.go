package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ExamRandomConfigModel = (*customExamRandomConfigModel)(nil)

type (
	// ExamRandomConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customExamRandomConfigModel.
	ExamRandomConfigModel interface {
		examRandomConfigModel
		withSession(session sqlx.Session) ExamRandomConfigModel
	}

	customExamRandomConfigModel struct {
		*defaultExamRandomConfigModel
	}
)

// NewExamRandomConfigModel returns a model for the database table.
func NewExamRandomConfigModel(conn sqlx.SqlConn) ExamRandomConfigModel {
	return &customExamRandomConfigModel{
		defaultExamRandomConfigModel: newExamRandomConfigModel(conn),
	}
}

func (m *customExamRandomConfigModel) withSession(session sqlx.Session) ExamRandomConfigModel {
	return NewExamRandomConfigModel(sqlx.NewSqlConnFromSession(session))
}
