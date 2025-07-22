package logic

import (
	"context"
	"exam-system/exam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AutoGradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAutoGradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AutoGradeLogic {
	return &AutoGradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AutoGradeLogic) AutoGrade() error {
	// todo: add your logic here and delete this line

	return nil
}
