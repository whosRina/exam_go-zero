package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListManualGradingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListManualGradingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListManualGradingLogic {
	return &ListManualGradingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListManualGradingLogic) ListManualGrading(req *types.ListManualGradingRequest, tokenString string) (*types.ListManualGradingResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 1 {
		return nil, errors.New("只有教师可以获取批阅列表")
	}

	// 查询考试记录，确保该考试由当前教师发布
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, req.ExamId)
	if err != nil {
		return nil, errors.New("考试不存在")
	}
	if exam.TeacherId != userId {
		return nil, errors.New("无权限获取该考试的批阅列表")
	}

	// 查询该考试下待批阅（状态为 "submitted"）的考试尝试记录
	attempts, err := l.svcCtx.ExamAttemptModel.FindSubmittedAttempts(l.ctx, req.ExamId)
	if err != nil {
		l.Logger.Errorf("查询批阅记录失败: %v", err)
		return nil, errors.New("获取批阅列表失败")
	}

	// 构造返回数据
	var attemptList []int64
	for _, at := range attempts {
		attemptList = append(attemptList, at.Id)
	}

	resp := &types.ListManualGradingResponse{
		AttemptList:  attemptList,
		PendingCount: len(attemptList),
	}
	return resp, nil
}
