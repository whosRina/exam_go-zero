package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/types"
	"time"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaperRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePaperRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePaperRuleLogic {
	return &UpdatePaperRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePaperRuleLogic) UpdatePaperRule(req *types.UpdatePaperRuleRequest, tokenString string) (*types.UpdatePaperRuleResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只有教师(userType == 1)才能更新规则
	if userType != 1 {
		return nil, errors.New("无权限更新规则")
	}

	// 查询规则是否存在
	rule, err := l.svcCtx.PaperRuleModel.FindOne(l.ctx, int64(req.RuleID))
	if err != nil {
		return nil, errors.New("规则不存在")
	}

	// 确保只有创建者可以修改
	if rule.CreatedBy != userId {
		return nil, errors.New("无权限修改他人创建的规则")
	}

	// 更新规则数据
	rule.Name = req.Name
	rule.TotalScore = int64(req.TotalScore)
	rule.BankId = req.BankID
	rule.NumQuestions = req.NumQuestions
	rule.ScoreConfig = req.ScoreConfig
	rule.UpdatedAt = time.Now()

	// 执行数据库更新
	err = l.svcCtx.PaperRuleModel.Update(l.ctx, rule)
	if err != nil {
		return nil, errors.New("更新规则失败")
	}

	return &types.UpdatePaperRuleResponse{
		Message: "规则更新成功",
	}, nil
}
