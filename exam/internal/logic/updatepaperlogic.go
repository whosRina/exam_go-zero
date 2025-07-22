package logic

import (
	"context"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/types"
	"time"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePaperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePaperLogic {
	return &UpdatePaperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePaperLogic) UpdatePaper(req *types.UpdatePaperRequest, tokenString string) (*types.UpdatePaperResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只有教师(userType == 1)才能更新试卷
	if userType != 1 {
		return nil, errors.New("无权限更新试卷")
	}

	// 查询试卷是否存在
	paper, err := l.svcCtx.PaperModel.FindOne(l.ctx, int64(req.PaperID))
	if err != nil {
		return nil, errors.New("试卷不存在")
	}

	// 确保只有试卷创建者可以修改
	if paper.CreatedBy != userId {
		return nil, errors.New("无权限修改他人创建的试卷")
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
			return nil, errors.New("题目分数必须大于0")
		}
		calculatedScore += q.Score
	}

	// 确保题目总分等于 TotalScore
	if calculatedScore != req.TotalScore {
		return nil, errors.New("题目分数总和必须等于试卷总分")
	}

	// 更新试卷数据
	paper.Name = req.Name
	paper.TotalScore = int64(req.TotalScore)
	paper.Questions = req.Questions
	paper.UpdatedAt = time.Now()

	// 执行数据库更新
	err = l.svcCtx.PaperModel.Update(l.ctx, paper)
	if err != nil {
		return nil, errors.New("更新试卷失败")
	}

	return &types.UpdatePaperResponse{
		Message: "试卷更新成功",
	}, nil
}
