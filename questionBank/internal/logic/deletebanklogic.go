package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/types"

	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBankLogic {
	return &DeleteBankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBankLogic) DeleteBank(req *types.DeleteBankRequest, tokenString string) (*types.DeleteBankResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 检查是否为教师
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限删除题库")
	}

	// 查询题库
	bank, err := l.svcCtx.QuestionBankModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		return nil, errors.New("查询题库失败")
	}

	// 如果是教师，检查是否为该教师创建的题库
	if bank.CreatedBy != userId {
		return nil, errors.New("您没有权限删除该题库")
	}

	// 删除题库中的所有问题
	err = l.svcCtx.QuestionBankModel.Delete(l.ctx, int64(req.Id))
	if err != nil {
		l.Logger.Errorf("删除题库问题失败:%v", err)
		return nil, errors.New("删除题库问题失败")
	}

	// 删除题库
	err = l.svcCtx.QuestionBankModel.Delete(l.ctx, int64(req.Id))
	if err != nil {
		l.Logger.Errorf("删除题库失败:%v", err)
		return nil, errors.New("删除题库失败")
	}

	// 返回成功响应
	return &types.DeleteBankResponse{
		Message: "题库删除成功",
	}, nil
}
