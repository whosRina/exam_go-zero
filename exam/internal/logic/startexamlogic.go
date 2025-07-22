package logic

import (
	"context"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"exam-system/exam/model"
	"math/rand"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartExamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartExamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartExamLogic {
	return &StartExamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartExamLogic) StartExam(req *types.StartExamRequest, tokenString string) (*types.StartExamResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	if userType != 2 {
		return nil, errors.New("只有学生可以开始考试")
	}

	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, req.ExamId)
	if err != nil {
		return nil, errors.New("考试不存在")
	}

	now := time.Now()
	if now.Before(exam.StartTime) {
		return nil, errors.New("考试还未开始")
	}
	if now.After(exam.EndTime) {
		return nil, errors.New("考试已结束")
	}

	existingAttempt, err := l.svcCtx.ExamAttemptModel.FindAttemptByExamAndStudent(l.ctx, req.ExamId, userId)
	if err != nil {
		return nil, errors.New("未找到考试记录")
	}

	if existingAttempt.Status != "not_started" && existingAttempt.Status != "ongoing" {
		return nil, errors.New("考试已完成")
	}

	var questions interface{}

	if exam.ExamType == "fixed" {
		paperId := exam.PaperId
		paper, err := l.svcCtx.PaperModel.FindOne(l.ctx, paperId)
		if err != nil {
			return nil, errors.New("关联的试卷不存在")
		}
		qErr := json.Unmarshal([]byte(paper.Questions), &questions)
		if qErr != nil {
			logx.Errorf("解析试卷题目失败: %v", qErr)
			return nil, errors.New("解析试卷题目失败")
		}
		existingAttempt.PaperId = paper.Id
		existingAttempt.GeneratedPaperId = -1

	} else if exam.ExamType == "random" {
		if existingAttempt.Status == "not_started" {
			// 只在第一次生成
			generatedPaper, err := GenerateRandomPaper(l.ctx, l.svcCtx, req.ExamId, userId)
			if err != nil {
				logx.Errorf("生成随机试卷失败: %v", err)
				return nil, errors.New("生成随机试卷失败")
			}

			qErr := json.Unmarshal([]byte(generatedPaper.Questions), &questions)
			if qErr != nil {
				logx.Errorf("解析随机试卷题目失败: %v", qErr)
				return nil, errors.New("解析试卷题目失败")
			}

			existingAttempt.GeneratedPaperId = generatedPaper.Id
			existingAttempt.PaperId = -1

		} else {
			// ongoing状态下读取之前生成的试卷
			generatedPaper, err := l.svcCtx.GeneratedPaperModel.FindOne(l.ctx, existingAttempt.GeneratedPaperId)
			if err != nil {
				return nil, errors.New("未找到已生成的试卷")
			}

			qErr := json.Unmarshal([]byte(generatedPaper.Questions), &questions)
			if qErr != nil {
				logx.Errorf("解析已生成试卷失败: %v", qErr)
				return nil, errors.New("解析试卷失败")
			}
		}
	}

	// 更新考试尝试记录状态为ongoing
	existingAttempt.Status = "ongoing"
	existingAttempt.StartTime = now
	err = l.svcCtx.ExamAttemptModel.Update(l.ctx, existingAttempt)
	if err != nil {
		logx.Errorf("更新考试尝试状态失败: %v", err)
		return nil, errors.New("更新考试状态失败")
	}

	return &types.StartExamResponse{
		Message:   "考试开始成功",
		AttemptId: existingAttempt.Id,
		Questions: questions,
	}, nil
}

func GenerateRandomPaper(ctx context.Context, svcCtx *svc.ServiceContext, examId, userId int64) (*model.GeneratedPaper, error) {
	// 查询考试信息
	exam, err := svcCtx.ExamModel.FindOne(ctx, examId)
	if err != nil {
		return nil, errors.New("考试信息查询失败")
	}

	// 查询随机规则
	rule, err := svcCtx.PaperRuleModel.FindOne(ctx, exam.PaperRuleId)
	if err != nil {
		return nil, errors.New("随机组卷规则未找到")
	}

	// 解析规则JSON
	numQuestions, err := parseRuleJSON(rule.NumQuestions)
	if err != nil {
		return nil, errors.New("解析题目数量规则失败")
	}
	scoreConfig, err := parseRuleJSON(rule.ScoreConfig)
	if err != nil {
		return nil, errors.New("解析题目分值规则失败")
	}

	// 按照规则随机抽取题目
	var selectedQuestions []map[string]int
	totalScore := 0
	rand.Seed(time.Now().UnixNano())

	for qType, count := range numQuestions {
		// 查询该题型的所有题目
		questions, err := svcCtx.QuestionModel.GetQuestionsByType(ctx, rule.BankId, qType)
		if err != nil {
			return nil, errors.New("查询题目失败")
		}

		// 随机打乱
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })

		// 选取前count道题目
		if len(questions) < count {
			return nil, errors.New("题目数量不足")
		}

		selected := questions[:count]
		for _, q := range selected {
			questionMap := map[string]int{
				"id":    int(q.Id),
				"score": scoreConfig[qType],
			}
			selectedQuestions = append(selectedQuestions, questionMap)
			totalScore += scoreConfig[qType]
		}
	}

	// 生成随机试卷JSON
	questionsJSON, err := json.Marshal(selectedQuestions)
	if err != nil {
		return nil, errors.New("生成试卷JSON失败")
	}

	// 6. 创建 generated_paper 记录
	paper := &model.GeneratedPaper{
		ExamId:     examId,
		StudentId:  userId,
		Questions:  string(questionsJSON),
		TotalScore: int64(totalScore),
		CreatedAt:  time.Now(),
	}
	result, err := svcCtx.GeneratedPaperModel.Insert(ctx, paper)
	if err != nil {
		return nil, errors.New("插入随机试卷失败")
	}

	// 获取插入ID
	paperId, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New("获取试卷ID失败")
	}
	paper.Id = paperId

	return paper, nil
}

// 解析JSON规则
func parseRuleJSON(jsonStr string) (map[int]int, error) {
	var data map[string]int
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	result := make(map[int]int)
	for key, value := range data {
		qType, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}
		result[qType] = value
	}

	return result, nil
}
