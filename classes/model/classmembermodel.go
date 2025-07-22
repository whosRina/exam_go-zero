package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ClassMemberModel = (*customClassMemberModel)(nil)

type (
	// ClassMemberModel is an interface to be customized, add more methods here,
	// and implement the added methods in customClassMemberModel.
	ClassMemberModel interface {
		classMemberModel
		withSession(session sqlx.Session) ClassMemberModel
	}

	customClassMemberModel struct {
		*defaultClassMemberModel
	}
)

// NewClassMemberModel returns a model for the database table.
func NewClassMemberModel(conn sqlx.SqlConn) ClassMemberModel {
	return &customClassMemberModel{
		defaultClassMemberModel: newClassMemberModel(conn),
	}
}

func (m *customClassMemberModel) withSession(session sqlx.Session) ClassMemberModel {
	return NewClassMemberModel(sqlx.NewSqlConnFromSession(session))
}
