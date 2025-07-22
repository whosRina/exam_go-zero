package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateExamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateExamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateExamLogic {
	return &UpdateExamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateExamLogic) UpdateExam(req *types.UpdateExamRequest, tokenString string) (*types.UpdateExamResponse, error) {
	// 身份认证,只允许教师身份更新考试
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil || userType != 1 {
		return nil, errors.New("教师身份认证失败")
	}

	// 参数校验,考试ID必须有效
	if req.ExamId <= 0 {
		return nil, errors.New("无效的考试ID")
	}
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

	// 验证该教师是否管理指定班级
	_, err = l.svcCtx.ClassModel.FindOneByTeacherId(l.ctx, req.ClassId, userId)
	if err != nil {
		return nil, errors.New("没有管理该班级的权限")
	}

	// 验证总分匹配（根据考试类型选择对应校验逻辑）
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

	// 获取原有考试记录，确保记录存在且归属当前教师
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, req.ExamId)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}
	if exam.TeacherId != userId {
		return nil, errors.New("无权限修改该考试")
	}

	// 更新考试记录字段
	exam.Name = req.Name
	exam.ClassId = req.ClassId
	exam.ExamType = req.ExamType
	exam.TotalScore = int64(req.TotalScore)
	exam.StartTime = startTime
	exam.EndTime = endTime
	exam.RequiresManualGrading = req.RequiresManualGrading
	exam.CanViewResults = req.CanViewResults
	if req.ExamType == "fixed" {
		exam.PaperId = req.PaperId
		exam.PaperRuleId = -1
	} else {
		exam.PaperRuleId = req.PaperRuleId
		exam.PaperId = -1
	}

	// 调用模型层更新记录
	err = l.svcCtx.ExamModel.Update(l.ctx, exam)
	if err != nil {
		return nil, errors.New("系统更新考试记录失败")
	}

	return &types.UpdateExamResponse{
		Message: "考试更新成功",
	}, nil
}
