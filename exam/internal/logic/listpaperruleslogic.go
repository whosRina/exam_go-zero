package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/types"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListPaperRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPaperRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPaperRulesLogic {
	return &ListPaperRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPaperRulesLogic) ListPaperRules(tokenString string) (*types.PaperRuleListResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的 JWT")
	}

	// 仅教师(userType == 1)可查看规则列表
	if userType != 1 {
		return nil, errors.New("无权限查看随机组卷规则")
	}

	// 查询该教师创建的规则列表
	rules, err := l.svcCtx.PaperRuleModel.FindAllByUserId(l.ctx, userId)
	if err != nil {
		return nil, errors.New("查询规则列表失败")
	}

	// 组装返回数据
	var ruleList []types.PaperRuleInfo
	for _, rule := range rules {
		ruleList = append(ruleList, types.PaperRuleInfo{
			Id:           int(rule.Id),
			Name:         rule.Name,
			TotalScore:   int(rule.TotalScore),
			CreatedBy:    int(rule.CreatedBy),
			BankId:       int(rule.BankId),
			NumQuestions: rule.NumQuestions,
			ScoreConfig:  rule.ScoreConfig,
			CreatedAt:    rule.CreatedAt,
			UpdatedAt:    rule.UpdatedAt,
		})
	}

	// 返回数据
	return &types.PaperRuleListResponse{
		RuleList: ruleList,
	}, nil
}
