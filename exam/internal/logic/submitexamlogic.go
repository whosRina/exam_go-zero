package logic

import (
	"context"
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"exam-system/exam/model"
	"fmt"
	"sort"
	"time"
	_ "time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitExamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitExamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitExamLogic {
	return &SubmitExamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SubmitExam 提交考试
func (l *SubmitExamLogic) SubmitExam(req *types.SubmitExamRequest, tokenString string) (*types.SubmitExamResponse, error) {
	// 解析JWT获取userId和userType
	now := time.Now()
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 2 {
		return nil, errors.New("只有学生可以提交考试")
	}

	// 查询考试尝试记录
	attempt, err := l.svcCtx.ExamAttemptModel.FindOne(l.ctx, req.AttemptId)
	if err != nil || attempt.StudentId != userId {
		return nil, errors.New("无权限提交该考试")
	}
	if attempt.Status == "submitted" || attempt.Status == "graded" {
		return nil, errors.New("不能重复考试")
	}

	var studentAnswers map[string]interface{}
	if err := json.Unmarshal([]byte(req.Answer), &studentAnswers); err != nil {
		return nil, errors.New("解析学生答案失败")
	}

	// 检查是否已有该考试的作答记录
	examAnswer, err := l.svcCtx.ExamAnswerModel.FindOneByAttempt(l.ctx, req.AttemptId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 没有找到记录，插入新答案
			answerJSON, _ := json.Marshal(studentAnswers)
			examAnswer = &model.ExamAnswer{
				AttemptId:     req.AttemptId,
				Answer:        string(answerJSON),
				GradingStatus: "pending",
				ScoreDetails:  "{}",
				SubmitTime:    now,
			}
			_, err = l.svcCtx.ExamAnswerModel.Insert(l.ctx, examAnswer)
			if err != nil {
				return nil, errors.New("提交答案失败")
			}
		} else {
			return nil, errors.New("查询答案记录失败")
		}
	} else {
		// 更新答案
		answerJSON, _ := json.Marshal(studentAnswers)
		examAnswer.Answer = string(answerJSON)
		examAnswer.SubmitTime = now
		examAnswer.GradingStatus = "pending"
	}

	// 评分逻辑
	scoreDetails, totalScore, err := l.GradeExam(attempt, examAnswer, studentAnswers)
	if err != nil {
		return nil, err
	}

	// 更新`score_details`
	scoreDetailsJSON, _ := json.Marshal(scoreDetails)
	examAnswer.ScoreDetails = string(scoreDetailsJSON)

	// 更新答案表
	if err := l.svcCtx.ExamAnswerModel.Update(l.ctx, examAnswer); err != nil {
		return nil, errors.New("更新答案失败")
	}

	// 更新考试尝试记录
	attempt.Score = int64(totalScore)
	attempt.SubmitTime = now
	if err := l.svcCtx.ExamAttemptModel.Update(l.ctx, attempt); err != nil {
		return nil, errors.New("提交考试失败")
	}

	return &types.SubmitExamResponse{
		Message: "考试提交成功",
	}, nil
}

// GradeExam 评分逻辑
func (l *SubmitExamLogic) GradeExam(attempt *model.ExamAttempt, examAnswer *model.ExamAnswer, studentAnswers map[string]interface{}) (map[int]int, int, error) {
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, attempt.ExamId)
	if err != nil {
		return nil, 0, errors.New("考试不存在")
	}

	// 评分状态更新
	if !exam.RequiresManualGrading {
		// 计算得分
		scoreDetails, totalScore, err := l.CalculateScore(attempt, studentAnswers)
		if err != nil {
			return nil, 0, err
		}
		attempt.Status = "graded"
		examAnswer.GradingStatus = "auto_scored"
		return scoreDetails, totalScore, nil
	} else {
		attempt.Status = "submitted"
		examAnswer.GradingStatus = "pending"
		return nil, 0, err
	}

}

// CalculateScore 计算考试分数
func (l *SubmitExamLogic) CalculateScore(attempt *model.ExamAttempt, studentAnswers map[string]interface{}) (map[int]int, int, error) {
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

	computedScore := 0
	scoreDetails := make(map[int]int) // 存储每道题的得分

	for _, entry := range questionEntries {
		question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, int64(entry.Id))
		if err != nil {
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

// isAnswerCorrect 根据题型比较学生答案与正确答案
func isAnswerCorrect(qType int, correct, student interface{}) bool {
	switch qType {
	case 1: // 单选题，正确答案为字符串
		return fmt.Sprintf("%v", student) == fmt.Sprintf("%v", correct)
	case 2: // 多选题，正确答案为数组
		// 强制检查是否为数组（如果是单个值，直接判错）
		_, isStudentArray := student.([]interface{}) // 如果student是JSON unmarshal后的数据，类型是 `[]interface{}`
		_, isCorrectArray := correct.([]interface{}) // 同理correct也必须是数组
		if !isStudentArray || !isCorrectArray {
			return false // 不是数组（比如传入 "A"），直接判 0 分
		}
		correctArr, ok1 := toStringSlice(correct)
		studentArr, ok2 := toStringSlice(student)
		if !ok1 || !ok2 {
			return false
		}
		// 排序后比较
		sort.Strings(correctArr)
		sort.Strings(studentArr)
		return fmt.Sprintf("%v", correctArr) == fmt.Sprintf("%v", studentArr)

	case 3: // 判断题，正确答案为布尔值
		correctBool := false
		studentBool := false

		// 处理正确答案
		switch v := correct.(type) {
		case bool:
			correctBool = v
		case string:
			correctBool = v == "true" || v == "A"
		}

		// 处理学生答案
		switch v := student.(type) {
		case bool:
			studentBool = v
		case string:
			studentBool = v == "true" || v == "A"
		default:
			// 避免意外类型直接判错
			return false
		}

		return correctBool == studentBool

	case 4:
		return fmt.Sprintf("%v", student) == fmt.Sprintf("%v", correct)
	default:
		return false
	}
}

// 尝试将接口转换为[]string
func toStringSlice(v interface{}) ([]string, bool) {
	switch val := v.(type) {
	case []interface{}:
		var result []string
		for _, item := range val {
			result = append(result, fmt.Sprintf("%v", item))
		}
		return result, true
	case []string:
		return val, true
	default:
		return nil, false
	}
}
