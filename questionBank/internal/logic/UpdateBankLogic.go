package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/svc"
	"exam-system/questionBank/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBankLogic {
	return &UpdateBankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBankLogic) UpdateBank(req *types.UpdateBankRequest, tokenString string) (*types.UpdateBankResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 检查是否为管理员或者教师
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限查看题库")
	}

	// 检查题库是否存在
	bank, err := l.svcCtx.QuestionBankModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		return nil, errors.New("查询题库失败")
	}

	// 检查题库是否属于该教师（通过createdBy字段）
	if bank.CreatedBy != userId {
		return nil, errors.New("没有权限更新该题库")
	}

	// 更新题库名称
	bank.Name = req.Name
	err = l.svcCtx.QuestionBankModel.Update(l.ctx, bank)
	if err != nil {
		l.Logger.Errorf("更新题库失败:%v", err)
		return nil, errors.New("更新题库失败")
	}

	// 返回成功响应
	return &types.UpdateBankResponse{
		Message: "题库更新成功",
	}, nil
}
