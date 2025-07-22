package logic

import (
	"context"
	"encoding/json"
	"errors"
	"exam-system/exam/model"

	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaperWithAnswersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaperWithAnswersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaperWithAnswersLogic {
	return &GetPaperWithAnswersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *GetPaperWithAnswersLogic) GetPaperWithAnswers(req *types.GetPaperWithAnswersRequest, tokenString string) (*types.PaperWithAnswersResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只有教师可以查看试卷详情
	if userType != 1 {
		return nil, errors.New("无权限查看试卷")
	}

	// 查询试卷并验证权限
	paper, err := l.svcCtx.PaperModel.FindPaperByUserId(l.ctx, userId, int64(req.PaperId))
	if err != nil {
		return nil, errors.New("试卷不存在或无权限查看此试卷")
	}

	// 解析试卷题目JSON
	var questionEntries []struct {
		QuestionId int `json:"id"`
		Score      int `json:"score"`
	}
	if err := json.Unmarshal([]byte(paper.Questions), &questionEntries); err != nil {
		l.Logger.Errorf("解析试卷题目失败: %v", err)
		return nil, errors.New("试卷题目解析失败")
	}

	// 查询题目信息并处理无效题目
	var validQuestions []types.QuestionWithAnswer
	var updatedQuestionEntries []struct {
		QuestionId int `json:"id"`
		Score      int `json:"score"`
	}
	var hasInvalidQuestions bool
	var newTotalScore int // 重新计算试卷总分

	for _, entry := range questionEntries {
		// 根据题目ID查询题目信息
		question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, int64(entry.QuestionId))
		if err != nil {
			// 题目不存在，标记需要更新试卷数据
			hasInvalidQuestions = true
			continue
		}

		// 题目存在，添加到新的列表
		validQuestions = append(validQuestions, types.QuestionWithAnswer{
			QuestionId: int(question.Id),
			Content:    question.Content,
			Type:       int(question.Type),
			Options:    question.Options,
			Answer:     question.Answer,
			Score:      entry.Score, // 题目分值
		})

		updatedQuestionEntries = append(updatedQuestionEntries, entry)

		// 重新计算试卷总分
		newTotalScore += entry.Score
	}

	// 如果试卷题目有变更，则更新paper表
	if hasInvalidQuestions {
		updatedQuestionsJSON, _ := json.Marshal(updatedQuestionEntries)

		// 构造更新数据
		updatedPaper := &model.Paper{
			Id:         paper.Id,
			Name:       paper.Name,
			TotalScore: int64(newTotalScore),
			Questions:  string(updatedQuestionsJSON),
		}

		// 调用PaperModel的Update方法
		err := l.svcCtx.PaperModel.Update(l.ctx, updatedPaper)
		if err != nil {
			l.Logger.Errorf("更新试卷题目或总分失败: %v", err)
		}
	}

	// 组装返回数据
	response := &types.PaperWithAnswersResponse{
		PaperId:    int(paper.Id),
		Name:       paper.Name,
		TotalScore: newTotalScore, // 返回最新的总分
		Questions:  validQuestions,
	}

	return response, nil
}
