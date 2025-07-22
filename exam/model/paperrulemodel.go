package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PaperRuleModel = (*customPaperRuleModel)(nil)

type (
	// PaperRuleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaperRuleModel.
	PaperRuleModel interface {
		paperRuleModel
		withSession(session sqlx.Session) PaperRuleModel
	}

	customPaperRuleModel struct {
		*defaultPaperRuleModel
	}
)

// NewPaperRuleModel returns a model for the database table.
func NewPaperRuleModel(conn sqlx.SqlConn) PaperRuleModel {
	return &customPaperRuleModel{
		defaultPaperRuleModel: newPaperRuleModel(conn),
	}
}

func (m *customPaperRuleModel) withSession(session sqlx.Session) PaperRuleModel {
	return NewPaperRuleModel(sqlx.NewSqlConnFromSession(session))
}
