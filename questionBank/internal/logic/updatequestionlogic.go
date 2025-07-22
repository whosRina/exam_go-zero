package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/types"

	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateQuestionLogic {
	return &UpdateQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateQuestionLogic) UpdateQuestion(req *types.UpdateQuestionRequest, tokenString string) (*types.UpdateQuestionResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 确保用户是教师才可更新题目
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限更新题目")
	}

	// 查询题目是否存在
	question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, int64(req.QuestionId))
	if err != nil || question == nil {
		return nil, errors.New("题目不存在")
	}
	// 检查题目是否属于该教师创建（通过createdBy字段）
	if question.CreatedBy != userId {
		return nil, errors.New("无权限更新该题目")
	}

	// 更新题目信息
	question.Content = req.Content
	question.Type = int64(req.Type)
	question.Options = req.Options
	question.Answer = req.Answer

	err = l.svcCtx.QuestionModel.Update(l.ctx, question)
	if err != nil {
		return nil, errors.New("更新题目失败")
	}

	return &types.UpdateQuestionResponse{
		Message: "题目更新成功",
	}, nil
}
