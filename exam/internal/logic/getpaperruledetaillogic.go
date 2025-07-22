package logic

import (
	"context"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaperRuleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaperRuleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaperRuleDetailLogic {
	return &GetPaperRuleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaperRuleDetailLogic) GetPaperRuleDetail() error {
	// todo: add your logic here and delete this line

	return nil
}
