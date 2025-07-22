package logic

import (
	"context"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"exam-system/exam/model"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitManualScoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitManualScoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitManualScoreLogic {
	return &SubmitManualScoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *SubmitManualScoreLogic) SubmitManualScore(req *types.SubmitManualScoreRequest, tokenString string) (*types.SubmitManualScoreResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 1 {
		return nil, errors.New("只有教师可以提交评分")
	}

	// 查询考试尝试记录
	attempt, err := l.svcCtx.ExamAttemptModel.FindOne(l.ctx, int64(req.AttemptId))
	if err != nil {
		return nil, errors.New("考试尝试记录不存在")
	}

	// 查询考试信息，确保教师有权限
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, attempt.ExamId)
	if err != nil || exam.TeacherId != userId {
		return nil, errors.New("无权限评分该考试")
	}

	// 获取学生答案记录
	examAnswer, err := l.svcCtx.ExamAnswerModel.FindOneByAttempt(l.ctx, int64(req.AttemptId))
	if err != nil {
		return nil, errors.New("未找到学生作答记录")
	}

	// 先执行自动评分，获取autoScoreDetails和autoTotalScore
	autoScoreDetails, _, err := l.CalculateScoreForNonManual(attempt, examAnswer)
	if err != nil {
		return nil, errors.New("自动评分失败")
	}

	// 解析已有的score_details（可能已有人工评分部分）
	var scoreDetails map[string]int
	if examAnswer.ScoreDetails != "" {
		if err := json.Unmarshal([]byte(examAnswer.ScoreDetails), &scoreDetails); err != nil {
			return nil, errors.New("解析已有的评分详情失败")
		}
		if scoreDetails == nil {
			scoreDetails = make(map[string]int) // 防止map是nil
		}
	} else {
		scoreDetails = make(map[string]int)
	}

	// 合并autoScoreDetails
	for qid, score := range autoScoreDetails {
		scoreDetails[fmt.Sprintf("%d", qid)] = score
	}

	if req.ManualScores != "{}" && req.ManualScores != "" {
		// 解析人工评分
		manualScores := make(map[string]int)
		if err := json.Unmarshal([]byte(req.ManualScores), &manualScores); err != nil {
			return nil, errors.New("解析人工评分失败: " + err.Error())
		}
		// 合并人工评分
		for qid, score := range manualScores {
			scoreDetails[qid] = score
		}
	}

	// 计算最终总分
	totalScore := 0
	for _, score := range scoreDetails {
		totalScore += score
	}

	// 更新examAnswer表
	scoreDetailsJSON, err := json.Marshal(scoreDetails)
	if err != nil {
		return nil, errors.New("JSON序列化失败")
	}
	examAnswer.ScoreDetails = string(scoreDetailsJSON)
	examAnswer.GradingStatus = "manual_scored"

	if err := l.svcCtx.ExamAnswerModel.Update(l.ctx, examAnswer); err != nil {
		return nil, errors.New("更新答案评分失败")
	}

	// 更新examAttempt表
	attempt.Score = int64(totalScore)
	attempt.Status = "graded"
	attempt.SubmitTime = time.Now()

	if err := l.svcCtx.ExamAttemptModel.Update(l.ctx, attempt); err != nil {
		return nil, errors.New("更新考试尝试失败")
	}

	// 返回成功响应
	return &types.SubmitManualScoreResponse{
		Message:    "评分提交成功",
		TotalScore: totalScore,
	}, nil
}

// CalculateScoreForNonManual 自动评分非人工阅卷的题目
func (l *SubmitManualScoreLogic) CalculateScoreForNonManual(attempt *model.ExamAttempt, examAnswer *model.ExamAnswer) (map[int]int, int, error) {
	var paperQuestionsJSON string
	if attempt.PaperId != -1 {
		paperRecord, err := l.svcCtx.PaperModel.FindOne(l.ctx, attempt.PaperId)
		if err != nil {
			return nil, 0, errors.New("关联试卷不存在")
		}
		paperQuestionsJSON = paperRecord.Questions
	} else if attempt.GeneratedPaperId != -1 {
		genPaper, err := l.svcCtx.GeneratedPaperModel.FindOne(l.ctx, attempt.GeneratedPaperId)
		if err != nil {
			return nil, 0, errors.New("随机试卷不存在")
		}
		paperQuestionsJSON = genPaper.Questions
	} else {
		return nil, 0, errors.New("试卷信息不完整")
	}

	var questionEntries []struct {
		Id    int `json:"id"`
		Score int `json:"score"`
	}
	if err := json.Unmarshal([]byte(paperQuestionsJSON), &questionEntries); err != nil {
		return nil, 0, errors.New("解析试卷题目失败")
	}

	var studentAnswers map[string]interface{}
	if err := json.Unmarshal([]byte(examAnswer.Answer), &studentAnswers); err != nil {
		return nil, 0, errors.New("解析学生答案失败")
	}

	computedScore := 0
	scoreDetails := make(map[int]int)

	for _, entry := range questionEntries {
		question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, int64(entry.Id))
		if err != nil || question.Type == 4 { // 过滤掉简答题
			continue
		}

		var correctAnswer interface{}
		if err := json.Unmarshal([]byte(question.Answer), &correctAnswer); err != nil {
			correctAnswer = question.Answer
		}

		studentAns := studentAnswers[fmt.Sprintf("%d", entry.Id)]
		if isAnswerCorrect(int(question.Type), correctAnswer, studentAns) {
			scoreDetails[entry.Id] = entry.Score
			computedScore += entry.Score
		} else {
			scoreDetails[entry.Id] = 0
		}
	}

	return scoreDetails, computedScore, nil
}
