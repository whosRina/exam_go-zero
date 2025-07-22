package logic

import (
	"context"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/types"
	"exam-system/exam/model"
	"time"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaperRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePaperRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaperRuleLogic {
	return &CreatePaperRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePaperRuleLogic) CreatePaperRule(req *types.CreatePaperRuleRequest, tokenString string) (*types.CreatePaperRuleResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只有教师才能创建规则
	if userType != 1 {
		return nil, errors.New("无权限创建规则")
	}

	// 校验总分数
	if req.TotalScore <= 0 {
		return nil, errors.New("试卷总分必须在大于0")
	}

	// 校验题型抽取数量(NumQuestions)
	var numQuestions map[string]int
	err = json.Unmarshal([]byte(req.NumQuestions), &numQuestions)
	if err != nil {
		return nil, errors.New("题型抽取数量格式错误")
	}

	// 校验题型分值(ScoreConfig)
	var scoreConfig map[string]int
	err = json.Unmarshal([]byte(req.ScoreConfig), &scoreConfig)
	if err != nil {
		return nil, errors.New("题型分值格式错误")
	}

	calculatedScore := 0
	for typeKey, score := range scoreConfig {
		count, exists := numQuestions[typeKey]
		if !exists || count < 0 {
			return nil, errors.New("存在无效的题型配置")
		}
		calculatedScore += count * score // 需要加上题目数量的乘积
	}

	if calculatedScore != req.TotalScore {
		return nil, errors.New("题型分值总和必须等于试卷总分")
	}

	// 确保题型分值总和等于试卷总分
	if calculatedScore != req.TotalScore {
		return nil, errors.New("题型分值总和必须等于试卷总分")
	}

	rule := &model.PaperRule{
		Name:         req.Name,
		TotalScore:   int64(req.TotalScore),
		CreatedBy:    userId,
		BankId:       int64(req.BankID),
		NumQuestions: req.NumQuestions, // 存储JSON格式的题型抽取数量
		ScoreConfig:  req.ScoreConfig,  // 存储JSON格式的题型分值
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 插入规则数据
	_, err = l.svcCtx.PaperRuleModel.Insert(l.ctx, rule)
	if err != nil {
		return nil, errors.New("创建随机组卷规则失败")
	}

	// 返回成功响应
	return &types.CreatePaperRuleResponse{
		Message: "随机组卷规则创建成功",
	}, nil
}
