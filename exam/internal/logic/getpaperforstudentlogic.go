package logic

import (
	"context"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/types"
	"strconv"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaperForStudentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaperForStudentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaperForStudentLogic {
	return &GetPaperForStudentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaperForStudentLogic) GetPaperForStudent(req *types.GetPaperForStudentRequest, tokenString string) (*types.PaperForStudentResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 2 {
		return nil, errors.New("无权限查看试卷")
	}

	// 查询考试尝试记录，确保记录存在且归属于当前学生
	attempt, err := l.svcCtx.ExamAttemptModel.FindOne(l.ctx, req.AttemptId)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}
	if attempt.StudentId != userId {
		return nil, errors.New("无权限查看该试卷")
	}
	if attempt.Status != "ongoing" {
		return nil, errors.New("考试已结束，无权限查看该试卷")
	}

	var (
		questions  []types.QuestionForStudent
		totalScore int64
		paperName  string
	)

	// 尝试获取学生已保存的作答记录
	var studentAnswers map[string]interface{}
	examAnswer, _ := l.svcCtx.ExamAnswerModel.FindOneByAttempt(l.ctx, req.AttemptId)
	if examAnswer != nil {
		if err := json.Unmarshal([]byte(examAnswer.Answer), &studentAnswers); err != nil {
			return nil, errors.New("解析学生答案失败")
		}
	}

	// 根据考试类型返回试卷及题目
	if attempt.PaperId != -1 {
		// 固定试卷
		paper, err := l.svcCtx.PaperModel.FindOne(l.ctx, attempt.PaperId)
		if err != nil {
			return nil, errors.New("试卷不存在")
		}
		paperName = paper.Name
		totalScore = paper.TotalScore

		var questionEntries []struct {
			QuestionId int `json:"id"`
			Score      int `json:"score"`
		}
		if err := json.Unmarshal([]byte(paper.Questions), &questionEntries); err != nil {
			l.Logger.Errorf("解析试卷题目失败: %v", err)
			return nil, errors.New("试卷题目解析失败")
		}
		for _, entry := range questionEntries {
			question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, int64(entry.QuestionId))
			if err != nil {
				continue
			}
			// 获取学生已提交的答案（若存在），否则为空字符串
			var ans interface{}
			if studentAnswers != nil {
				ans = studentAnswers[strconv.Itoa(entry.QuestionId)]
			}
			questions = append(questions, types.QuestionForStudent{
				QuestionId: int(question.Id),
				Content:    question.Content,
				Type:       int(question.Type),
				Options:    question.Options,
				Score:      entry.Score,
				Answer:     ans, // 返回学生答案
			})
		}
	} else if attempt.GeneratedPaperId != -1 {
		// 随机试卷
		genPaper, err := l.svcCtx.GeneratedPaperModel.FindOne(l.ctx, attempt.GeneratedPaperId)
		if err != nil {
			return nil, errors.New("随机试卷不存在")
		}
		totalScore = genPaper.TotalScore
		paperName = "随机试卷"
		var questionEntries []struct {
			QuestionId int `json:"id"`
			Score      int `json:"score"`
		}
		if err := json.Unmarshal([]byte(genPaper.Questions), &questionEntries); err != nil {
			l.Logger.Errorf("解析随机试卷题目失败: %v", err)
			return nil, errors.New("解析随机试卷题目失败")
		}
		for _, entry := range questionEntries {
			question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, int64(entry.QuestionId))
			if err != nil {
				continue
			}
			key := strconv.Itoa(entry.QuestionId)
			var ans interface{}
			if studentAnswers != nil {
				ans = studentAnswers[key]
			}
			questions = append(questions, types.QuestionForStudent{
				QuestionId: int(question.Id),
				Content:    question.Content,
				Type:       int(question.Type),
				Options:    question.Options,
				Score:      entry.Score,
				Answer:     ans,
			})
		}
	} else {
		return nil, errors.New("试卷信息不完整")
	}
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, attempt.ExamId)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}

	response := &types.PaperForStudentResponse{
		StartTime:  exam.StartTime,
		EndTime:    exam.EndTime,
		PaperName:  paperName,
		ExamId:     int(exam.Id),
		TotalScore: int(totalScore),
		Questions:  questions,
	}
	return response, nil
}
