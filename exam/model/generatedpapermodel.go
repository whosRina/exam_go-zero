package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GeneratedPaperModel = (*customGeneratedPaperModel)(nil)

type (
	// GeneratedPaperModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGeneratedPaperModel.
	GeneratedPaperModel interface {
		generatedPaperModel
		withSession(session sqlx.Session) GeneratedPaperModel
	}

	customGeneratedPaperModel struct {
		*defaultGeneratedPaperModel
	}
)

// NewGeneratedPaperModel returns a model for the database table.
func NewGeneratedPaperModel(conn sqlx.SqlConn) GeneratedPaperModel {
	return &customGeneratedPaperModel{
		defaultGeneratedPaperModel: newGeneratedPaperModel(conn),
	}
}

func (m *customGeneratedPaperModel) withSession(session sqlx.Session) GeneratedPaperModel {
	return NewGeneratedPaperModel(sqlx.NewSqlConnFromSession(session))
}
