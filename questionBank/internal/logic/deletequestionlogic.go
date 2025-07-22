package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/svc"
	"exam-system/questionBank/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionLogic {
	return &DeleteQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteQuestionLogic) DeleteQuestion(req *types.DeleteQuestionRequest, tokenString string) (*types.DeleteQuestionResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 根据用户类型鉴别用户是否有权限删除题目
	if userType != 1 { // 只有教师才有权限删除题目
		return nil, errors.New("无权限删除题目")
	}

	// 检查题目是否存在
	question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, req.QuestionId)
	if err != nil {
		return nil, errors.New("查询题目失败")
	}

	// 确保查询到的题目属于该教师
	if question.CreatedBy != userId {
		return nil, errors.New("无权限删除此题目")
	}

	// 删除题目
	err = l.svcCtx.QuestionModel.Delete(l.ctx, req.QuestionId)
	if err != nil {
		return nil, errors.New("删除题目失败")
	}

	// 返回删除成功的响应
	return &types.DeleteQuestionResponse{
		Message: "题目删除成功",
	}, nil
}
