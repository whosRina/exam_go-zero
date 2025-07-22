package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"exam-system/exam/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitAnswerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitAnswerLogic {
	return &SubmitAnswerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitAnswerLogic) SubmitAnswer(req *types.SubmitAnswerRequest, tokenString string) (*types.SubmitAnswerResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 2 {
		return nil, errors.New("只有学生可以提交答案")
	}

	// 查询考试尝试记录，确保记录存在且归属于当前学生
	attempt, err := l.svcCtx.ExamAttemptModel.FindOne(l.ctx, req.AttemptId)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}
	if attempt.StudentId != userId {
		return nil, errors.New("无权限提交该考试答案")
	}

	// 查找是否已有对应的作答记录
	// 查找作答记录
	examAnswer, err := l.svcCtx.ExamAnswerModel.FindOneByAttempt(l.ctx, req.AttemptId)
	if err != nil {
		l.Logger.Errorf("未找到作答记录，将创建新记录:%v", err)
		// 创建新作答记录
		newRecord := &model.ExamAnswer{
			AttemptId:     req.AttemptId,
			Answer:        req.Answer,
			GradingStatus: "pending",
			ScoreDetails:  "{}",
			SubmitTime:    time.Now(),
		}
		_, err = l.svcCtx.ExamAnswerModel.Insert(l.ctx, newRecord)
		if err != nil {
			l.Logger.Errorf("插入新答案失败:%v", err)
			return nil, errors.New("提交答案失败")
		}
	} else {
		// 更新已有的作答记录
		examAnswer.Answer = req.Answer
		examAnswer.SubmitTime = time.Now()
		err = l.svcCtx.ExamAnswerModel.Update(l.ctx, examAnswer)
		if err != nil {
			l.Logger.Errorf("更新答案失败: %v", err)
			return nil, errors.New("更新答案失败")
		}
	}

	l.Logger.Infof("学生%d成功提交attemptId:%d的答案", userId, req.AttemptId)
	return &types.SubmitAnswerResponse{
		Message: "答案提交成功",
	}, nil
}
