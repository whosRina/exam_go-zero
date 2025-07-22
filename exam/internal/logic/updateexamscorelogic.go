package logic

import (
	"context"
	"exam-system/exam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateExamScoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateExamScoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateExamScoreLogic {
	return &UpdateExamScoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateExamScoreLogic) UpdateExamScore() error {
	// todo: add your logic here and delete this line

	return nil
}
