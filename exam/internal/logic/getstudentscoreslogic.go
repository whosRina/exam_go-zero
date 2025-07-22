package logic

import (
	"context"
	"exam-system/exam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStudentScoresLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStudentScoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStudentScoresLogic {
	return &GetStudentScoresLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStudentScoresLogic) GetStudentScores() error {
	// todo: add your logic here and delete this line

	return nil
}
