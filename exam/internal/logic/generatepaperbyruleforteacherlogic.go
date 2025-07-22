package logic

import (
	"context"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GeneratePaperByRuleForTeacherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeneratePaperByRuleForTeacherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeneratePaperByRuleForTeacherLogic {
	return &GeneratePaperByRuleForTeacherLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeneratePaperByRuleForTeacherLogic) GeneratePaperByRuleForTeacher() error {
	// todo: add your logic here and delete this line

	return nil
}
