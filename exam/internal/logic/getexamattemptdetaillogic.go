package logic

import (
	"context"
	"exam-system/exam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExamAttemptDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExamAttemptDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExamAttemptDetailLogic {
	return &GetExamAttemptDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExamAttemptDetailLogic) GetExamAttemptDetail() error {
	// todo: add your logic here and delete this line

	return nil
}
