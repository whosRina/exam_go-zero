package logic

import (
	"context"
	"exam-system/exam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetExamQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetExamQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetExamQuestionsLogic {
	return &SetExamQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetExamQuestionsLogic) SetExamQuestions() error {
	// todo: add your logic here and delete this line

	return nil
}
