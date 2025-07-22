package logic

import (
	"context"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaperDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaperDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaperDetailLogic {
	return &GetPaperDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaperDetailLogic) GetPaperDetail() error {
	// todo: add your logic here and delete this line

	return nil
}
