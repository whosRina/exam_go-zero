package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/types"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePaperRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePaperRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePaperRuleLogic {
	return &DeletePaperRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePaperRuleLogic) DeletePaperRule(req *types.DeletePaperRuleRequest, tokenString string) (*types.DeletePaperRuleResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只有教师才能删除随机规则
	if userType != 1 {
		return nil, errors.New("无权限删除随机组卷规则")
	}

	// 查询规则是否存在
	rule, err := l.svcCtx.PaperRuleModel.FindOne(l.ctx, int64(req.RuleID))
	if err != nil {
		return nil, errors.New("规则不存在")
	}

	// 确保只有规则创建者可以删除
	if rule.CreatedBy != userId {
		return nil, errors.New("无权限删除他人创建的规则")
	}

	// 执行删除操作
	err = l.svcCtx.PaperRuleModel.Delete(l.ctx, int64(req.RuleID))
	if err != nil {
		return nil, errors.New("删除规则失败")
	}

	return &types.DeletePaperRuleResponse{
		Message: "规则删除成功",
	}, nil
}
