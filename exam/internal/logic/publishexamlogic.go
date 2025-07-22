package logic

import (
	"context"
	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishExamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishExamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishExamLogic {
	return &PublishExamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishExamLogic) PublishExam() error {
	// todo: add your logic here and delete this line

	return nil
}
