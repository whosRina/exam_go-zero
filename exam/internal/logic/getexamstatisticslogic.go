package logic

import (
	"context"
	"exam-system/exam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExamStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExamStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExamStatisticsLogic {
	return &GetExamStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExamStatisticsLogic) GetExamStatistics() error {
	// todo: add your logic here and delete this line

	return nil
}
