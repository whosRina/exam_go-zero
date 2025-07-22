package logic

import (
	"context"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManualGradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManualGradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManualGradeLogic {
	return &ManualGradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManualGradeLogic) ManualGrade(req *types.ManualGradeRequest, tokenString string) (*types.ManualGradeResponse, error) {
	// 解析 JWT，确保调用者为教师
	teacherId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 1 {
		return nil, errors.New("只有教师可以获取人工批改详情")
	}

	// 查询考试尝试记录
	attempt, err := l.svcCtx.ExamAttemptModel.FindOne(l.ctx, req.AttemptId)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}

	// 验证该考试是否由当前教师发布
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, attempt.ExamId)
	if err != nil {
		return nil, errors.New("考试不存在")
	}
	if exam.TeacherId != teacherId {
		return nil, errors.New("无权限获取该考试的批改详情")
	}

	// 查询考试作答记录（通过 attempt_id 获取答案记录）
	examAnswer, err := l.svcCtx.ExamAnswerModel.FindOneByAttempt(l.ctx, req.AttemptId)
	if err != nil {
		return nil, errors.New("考试作答记录不存在")
	}

	// 获取试卷题目 JSON（固定试卷或随机试卷）
	var paperQuestionsJSON string
	if attempt.PaperId != -1 {
		paperRecord, err := l.svcCtx.PaperModel.FindOne(l.ctx, attempt.PaperId)
		if err != nil {
			return nil, errors.New("关联试卷不存在")
		}
		paperQuestionsJSON = paperRecord.Questions
	} else if attempt.GeneratedPaperId != -1 {
		genPaper, err := l.svcCtx.GeneratedPaperModel.FindOne(l.ctx, attempt.GeneratedPaperId)
		if err != nil {
			return nil, errors.New("随机试卷不存在")
		}
		paperQuestionsJSON = genPaper.Questions
	} else {
		return nil, errors.New("试卷信息不完整")
	}

	// 解析试卷题目，期望格式为：[{\"id\": 19, \"score\": 25}, ...]
	var questionEntries []struct {
		Id    int `json:"id"`
		Score int `json:"score"`
	}
	if err := json.Unmarshal([]byte(paperQuestionsJSON), &questionEntries); err != nil {
		return nil, errors.New("解析试卷题目失败")
	}

	// 解析学生答案，格式为 JSON 对象：{\"19\":\"C\", ...}
	var studentAnswers map[string]interface{}
	if err := json.Unmarshal([]byte(examAnswer.Answer), &studentAnswers); err != nil {
		studentAnswers = make(map[string]interface{})
	}

	// 遍历题目列表，筛选出简答题 (type == 4) 并组装返回数据
	var details []types.ShortAnswerDetail
	for _, entry := range questionEntries {
		question, err := l.svcCtx.QuestionModel.FindOne(l.ctx, int64(entry.Id))
		if err != nil {
			continue
		}
		// 仅处理简答题，简答题需要人工批改（type == 4）
		if question.Type != 4 {
			continue
		}
		// 获取学生答案（键为题目ID的字符串形式）
		studentAns := ""
		if ans, ok := studentAnswers[fmt.Sprintf("%d", entry.Id)]; ok {
			studentAns = fmt.Sprintf("%v", ans)
		}
		// 组装返回的题目详情数据
		detail := types.ShortAnswerDetail{
			QuestionId:    entry.Id,
			Content:       question.Content,
			Answer:        question.Answer,
			Score:         entry.Score,
			StudentAnswer: studentAns,
		}
		details = append(details, detail)
	}

	resp := &types.ManualGradeResponse{
		Questions: details,
	}
	return resp, nil

}
