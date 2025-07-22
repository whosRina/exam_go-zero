package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/types"
	"exam-system/questionBank/model"
	"time"

	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBankLogic {
	return &CreateBankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBankLogic) CreateBank(req *types.CreateBankRequest, tokenString string) (*types.CreateBankResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 确保用户是教师才可以创建题库
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限创建题库")
	}

	// 创建题库
	bank := &model.QuestionBank{
		Name:      req.BankName,
		CreatedBy: userId,
		CreatedAt: time.Now(), // 当前时间作为创建时间
	}

	_, err = l.svcCtx.QuestionBankModel.Insert(l.ctx, bank)
	if err != nil {
		return nil, errors.New("创建题库失败")
	}

	return &types.CreateBankResponse{
		Message: "题库创建成功",
	}, nil
}
