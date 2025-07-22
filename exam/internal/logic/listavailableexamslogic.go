package logic

import (
	"context"
	"exam-system/exam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAvailableExamsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAvailableExamsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAvailableExamsLogic {
	return &ListAvailableExamsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAvailableExamsLogic) ListAvailableExams() error {
	// todo: add your logic here and delete this line

	return nil
}
