package logic

import (
	"context"

	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuestionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionDetailLogic {
	return &QuestionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuestionDetailLogic) QuestionDetail() error {
	// todo: add your logic here and delete this line

	return nil
}
