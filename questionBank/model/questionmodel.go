package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ QuestionModel = (*customQuestionModel)(nil)

type (
	// QuestionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionModel.
	QuestionModel interface {
		questionModel
		withSession(session sqlx.Session) QuestionModel
	}

	customQuestionModel struct {
		*defaultQuestionModel
	}
)

// NewQuestionModel returns a model for the database table.
func NewQuestionModel(conn sqlx.SqlConn) QuestionModel {
	return &customQuestionModel{
		defaultQuestionModel: newQuestionModel(conn),
	}
}

func (m *customQuestionModel) withSession(session sqlx.Session) QuestionModel {
	return NewQuestionModel(sqlx.NewSqlConnFromSession(session))
}
