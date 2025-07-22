package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteExamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteExamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteExamLogic {
	return &DeleteExamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteExamLogic) DeleteExam(req *types.DeleteExamRequest, tokenString string) (*types.DeleteExamResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只有教师才能删除考试
	if userType != 1 {
		return nil, errors.New("无权限删除考试")
	}

	// 查询考试是否存在
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, req.ExamId)
	if err != nil {
		return nil, errors.New("考试不存在")
	}

	// 确保只有考试创建者可以删除
	if exam.TeacherId != userId {
		return nil, errors.New("无权限删除该考试")
	}

	// 执行删除操作
	err = l.svcCtx.ExamModel.Delete(l.ctx, req.ExamId)
	if err != nil {
		return nil, errors.New("删除考试失败")
	}

	return &types.DeleteExamResponse{
		Message: "考试删除成功",
	}, nil
}
