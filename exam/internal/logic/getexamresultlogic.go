package logic

import (
	"context"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExamResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExamResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExamResultLogic {
	return &GetExamResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExamResultLogic) GetExamResult(req *types.GetExamResultRequest, tokenString string) (*types.ExamResultResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 2 {
		return nil, errors.New("无权限查看考试结果")
	}

	// 查询考试尝试记录
	attempt, err := l.svcCtx.ExamAttemptModel.FindAttemptByExamAndStudent(l.ctx, req.ExamId, userId)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}
	if attempt.Status != "graded" && attempt.Status != "submitted" {
		return nil, errors.New("考试未结束，无法查看成绩")
	}

	// 获取考试信息
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, attempt.ExamId)
	if err != nil {
		return nil, errors.New("考试信息不存在")
	}

	// 获取学生答案和评分详情
	examAnswer, err := l.svcCtx.ExamAnswerModel.FindOneByAttempt(l.ctx, attempt.Id)
	if err != nil {
		return nil, errors.New("学生答案记录不存在")
	}

	var studentAnswers map[string]interface{}
	var studentScores map[string]int
	if err := json.Unmarshal([]byte(examAnswer.Answer), &studentAnswers); err != nil {
		return nil, errors.New("解析学生答案失败")
	}
	if err := json.Unmarshal([]byte(examAnswer.ScoreDetails), &studentScores); err != nil {
		return nil, errors.New("解析评分详情失败")
	}

	var (
		questions  []types.QuestionResult
		totalScore int64
	)

	// 获取试卷信息（固定试卷或随机试卷）
	var questionEntries []struct {
		QuestionId int `json:"id"`
		Score      int `json:"score"`
	}

	if attempt.PaperId != -1 {
		// 固定试卷
		paper, err := l.svcCtx.PaperModel.FindOne(l.ctx, attempt.PaperId)
		if err != nil {
			return nil, errors.New("试卷不存在")
		}

		totalScore = paper.TotalScore

		if err := json.Unmarshal([]byte(paper.Questions), &questionEntries); err != nil {
			return nil, errors.New("解析试卷题目失败")
		}
	} else if attempt.GeneratedPaperId != -1 {
		// 随机试卷
		genPaper, err := l.svcCtx.GeneratedPaperModel.FindOne(l.ctx, attempt.GeneratedPaperId)
		if err != nil {
			return nil, errors.New("随机试卷不存在")
		}

		totalScore = genPaper.TotalScore

		if err := json.Unmarshal([]byte(genPaper.Questions), &questionEntries); err != nil {
			return nil, errors.New("解析随机试卷题目失败")
		}
	} else {
		return nil, errors.New("试卷信息不完整")
	}

	// 计算评分状态
	gradingStatus := "pending" // 默认未评分
	if attempt.Status == "graded" {
		gradingStatus = "graded"
	}

	// 如果考试不允许查看试题信息questions置空
	if exam.CanViewResults && attempt.Status == "graded" {
		// 组装返回的考试题目详情
		for _, entry := range questionEntries {
			question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, int64(entry.QuestionId))
			if err != nil {
				continue
			}

			// 获取学生答案和得分
			questionIdStr := strconv.Itoa(entry.QuestionId)
			var studentAnswer interface{}
			var studentScore int

			if studentAnswers != nil {
				studentAnswer = studentAnswers[questionIdStr]
			}
			if studentScores != nil {
				studentScore = studentScores[questionIdStr]
			}
			var formattedAnswer string
			if err := json.Unmarshal([]byte(question.Answer), &formattedAnswer); err != nil {
				formattedAnswer = question.Answer // 解析失败时，直接返回原始值
			}

			questions = append(questions, types.QuestionResult{
				QuestionId:    int(question.Id),
				Content:       question.Content,
				Type:          int(question.Type),
				Options:       question.Options,
				TotalScore:    entry.Score,
				StudentAnswer: studentAnswer,
				Answer:        formattedAnswer,
				StudentScore:  studentScore,
			})
		}
	}

	// 返回考试结果
	response := &types.ExamResultResponse{
		ExamId:         attempt.ExamId,
		ExamName:       exam.Name,
		TotalScore:     int(totalScore),
		ExamScore:      int(attempt.Score),
		CanViewResults: exam.CanViewResults,
		GradingStatus:  gradingStatus, // 评分状态
		Questions:      questions,     // 题目信息
	}
	return response, nil
}
