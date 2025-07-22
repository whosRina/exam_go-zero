package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"exam-system/exam/model"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateExamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateExamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateExamLogic {
	return &CreateExamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateExamLogic) CreateExam(req *types.CreateExamRequest, tokenString string) (*types.CreateExamResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil || userType != 1 {
		return nil, errors.New("教师身份认证失败")
	}

	// 基础参数校验
	if req.Name == "" {
		return nil, errors.New("考试名称不能为空")
	}
	if req.TotalScore <= 0 {
		return nil, errors.New("考试总分需大于0")
	}
	if req.ExamType != "fixed" && req.ExamType != "random" {
		return nil, errors.New("考试类型必须为fixed或random")
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime, err := time.ParseInLocation("2006-01-02 15:04", req.StartTime, loc)
	if err != nil {
		return nil, errors.New("开始时间格式错误，要求格式：2006-01-02 15:04")
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04", req.EndTime, loc)
	if err != nil {
		return nil, errors.New("结束时间格式错误，要求格式：2006-01-02 15:04")
	}
	// 验证班级归属
	_, err = l.svcCtx.ClassModel.FindOneByTeacherId(l.ctx, req.ClassId, userId)
	if err != nil {
		return nil, errors.New("没有管理该班级的权限")
	}

	// 总分匹配验证
	switch req.ExamType {
	case "fixed":
		paper, err := l.svcCtx.PaperModel.FindOne(l.ctx, req.PaperId)
		if err != nil {
			return nil, errors.New("关联试卷不存在")
		}
		if paper.TotalScore != int64(req.TotalScore) {
			return nil, fmt.Errorf("试卷总分不符(试卷:%d 考试:%d)", paper.TotalScore, req.TotalScore)
		}

	case "random":
		rule, err := l.svcCtx.PaperRuleModel.FindOne(l.ctx, req.PaperRuleId)
		if err != nil {
			return nil, errors.New("组卷规则不存在")
		}
		if rule.TotalScore != int64(req.TotalScore) {
			return nil, fmt.Errorf("规则总分不符(规则:%d 考试:%d)", rule.TotalScore, req.TotalScore)
		}
	}

	// 创建考试记录
	exam := &model.Exam{
		Name:                  req.Name,
		TeacherId:             userId,
		ClassId:               req.ClassId,
		ExamType:              req.ExamType,
		TotalScore:            int64(req.TotalScore),
		StartTime:             startTime,
		EndTime:               endTime,
		RequiresManualGrading: req.RequiresManualGrading,
		CanViewResults:        req.CanViewResults,
	}
	if req.ExamType == "fixed" {
		exam.PaperId = req.PaperId
		exam.PaperRuleId = -1
	} else {
		exam.PaperRuleId = req.PaperRuleId
		exam.PaperId = -1
	}

	// 插入考试记录
	ret, err := l.svcCtx.ExamModel.Insert(l.ctx, exam)
	if err != nil {
		return nil, errors.New("系统保存考试记录失败")
	}

	// 获取新创建的考试 ID
	examId, err := ret.LastInsertId()
	if err != nil {
		return nil, errors.New("获取考试ID失败")
	}

	// 查询班级成员
	members, err := l.svcCtx.ClassMemberModel.FindByClassId(l.ctx, req.ClassId)
	if err != nil {
		return nil, errors.New("获取班级成员失败")
	}

	// 为每个学生创建考试尝试记录
	for _, member := range members {
		attempt := &model.ExamAttempt{
			ExamId:     examId, // 这里插入 examId
			StudentId:  member.UserId,
			Status:     "not_started", // 新增未开始状态
			CreatedAt:  time.Now(),
			StartTime:  time.Unix(0, 0).In(time.FixedZone("CST", 8*3600)), // 1970-01-01 08:00:00
			SubmitTime: time.Unix(0, 0).In(time.FixedZone("CST", 8*3600)),
		}

		_, err := l.svcCtx.ExamAttemptModel.Insert(l.ctx, attempt)
		if err != nil {
			logx.Errorf("插入考试尝试记录失败: %v", err)
		}
	}

	return &types.CreateExamResponse{
		Message: "考试创建成功",
	}, nil

}
