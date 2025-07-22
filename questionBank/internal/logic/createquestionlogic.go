package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/types"
	"exam-system/questionBank/model"

	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionLogic {
	return &CreateQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateQuestionLogic) CreateQuestion(req *types.CreateQuestionRequest, tokenString string) (*types.CreateQuestionResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 确保用户是教师才可创建题目
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限创建题目")
	}

	// 检查题库是否存在
	bank, err := l.svcCtx.QuestionBankModel.FindById(l.ctx, req.BankId)
	if err != nil || bank == nil {
		return nil, errors.New("题库不存在")
	}

	question := &model.Question{
		BankId:    int64(req.BankId),
		Content:   req.Content,
		Type:      int64(req.Type),
		Options:   req.Options, // 默认值
		Answer:    req.Answer,
		CreatedBy: userId,
	}

	_, err = l.svcCtx.QuestionModel.Insert(l.ctx, question)
	if err != nil {
		return nil, errors.New("创建题目失败")
	}

	return &types.CreateQuestionResponse{
		Message: "题目创建成功",
	}, nil
}
