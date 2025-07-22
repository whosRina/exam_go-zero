package logic

import (
	"context"
	"exam-system/exam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetExamRandomConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetExamRandomConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetExamRandomConfigLogic {
	return &SetExamRandomConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetExamRandomConfigLogic) SetExamRandomConfig() error {
	// todo: add your logic here and delete this line

	return nil
}
