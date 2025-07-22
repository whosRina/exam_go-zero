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

type CreatePaperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePaperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaperLogic {
	return &CreatePaperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePaperLogic) CreatePaper(req *types.CreatePaperRequest, tokenString string) (*types.CreatePaperResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只有教师才能创建试卷
	if userType != 1 {
		return nil, errors.New("无权限创建试卷")
	}
	// 校验总分数
	if req.TotalScore <= 0 {
		return nil, errors.New("试卷总分必须大于0")
	}

	// 解析questions JSON，检查题目总分是否匹配TotalScore
	var questions []struct {
		Id    int `json:"id"`
		Score int `json:"score"`
	}
	err = json.Unmarshal([]byte(req.Questions), &questions)
	if err != nil {
		return nil, errors.New("题目列表格式错误")
	}

	// 计算题目总分
	var calculatedScore int
	for _, q := range questions {
		if q.Score <= 0 {
			return nil, errors.New("题目分数必须大于 0")
		}
		calculatedScore += q.Score
	}

	// 确保题目总分等于TotalScore
	if calculatedScore != req.TotalScore {
		return nil, errors.New("题目分数总和必须等于试卷总分")
	}

	// 组装试卷数据
	paper := &model.Paper{
		Name:       req.Name,
		TotalScore: int64(req.TotalScore),
		CreatedBy:  userId,
		Questions:  req.Questions, // 存储JSON格式的题目列表
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 插入试卷数据
	_, err = l.svcCtx.PaperModel.Insert(l.ctx, paper)
	if err != nil {
		return nil, errors.New("创建试卷失败")
	}

	// 返回成功响应
	return &types.CreatePaperResponse{
		Message: "试卷创建成功",
	}, nil
}
