package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ QuestionBankModel = (*customQuestionBankModel)(nil)

type (
	// QuestionBankModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionBankModel.
	QuestionBankModel interface {
		questionBankModel
		withSession(session sqlx.Session) QuestionBankModel
	}

	customQuestionBankModel struct {
		*defaultQuestionBankModel
	}
)

// NewQuestionBankModel returns a model for the database table.
func NewQuestionBankModel(conn sqlx.SqlConn) QuestionBankModel {
	return &customQuestionBankModel{
		defaultQuestionBankModel: newQuestionBankModel(conn),
	}
}

func (m *customQuestionBankModel) withSession(session sqlx.Session) QuestionBankModel {
	return NewQuestionBankModel(sqlx.NewSqlConnFromSession(session))
}
