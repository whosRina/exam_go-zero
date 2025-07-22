package logic

import (
	"context"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaperByRuleForStudentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaperByRuleForStudentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaperByRuleForStudentLogic {
	return &GetPaperByRuleForStudentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaperByRuleForStudentLogic) GetPaperByRuleForStudent() error {
	// todo: add your logic here and delete this line

	return nil
}
