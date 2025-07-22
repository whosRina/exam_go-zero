package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExamDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExamDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExamDetailLogic {
	return &ExamDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExamDetailLogic) ExamDetail(req *types.ExamDetailRequest, tokenString string) (*types.ExamDetailResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 仅允许教师访问
	if userType != 1 {
		return nil, errors.New("无权限查看该考试详情")
	}

	// 查询考试信息
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, req.ExamId)
	if err != nil {
		return nil, errors.New("考试不存在")
	}

	// 验证教师是否为考试的创建者
	if exam.TeacherId != userId {
		return nil, errors.New("你不是该考试的创建者，无法查看详情")
	}

	// 查询班级名称
	className := ""
	class, err := l.svcCtx.ClassModel.FindOne(l.ctx, exam.ClassId)
	if err == nil {
		className = class.Name
	} else {
		logx.Errorf("查询班级名称失败: %v", err)
	}

	// 查询创建者名称（教师名称）
	createBy := ""
	teacher, err := l.svcCtx.UsersModel.FindOne(l.ctx, exam.TeacherId)
	if err == nil {
		createBy = teacher.Name
	} else {
		logx.Errorf("查询创建者名称失败: %v", err)
	}

	// 格式化时间，使用 "2006-01-02 15:04" 格式
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTimeStr := exam.StartTime.In(loc).Format("2006-01-02 15:04")
	endTimeStr := exam.EndTime.In(loc).Format("2006-01-02 15:04")

	// 查询该考试下所有考生的考试尝试记录
	attempts, err := l.svcCtx.ExamAttemptModel.FindAttemptByExam(l.ctx, exam.Id)
	if err != nil {
		return nil, errors.New("查询考试尝试记录失败")
	}

	// 组装每个考生的考试状态和成绩信息
	var studentResults []types.ExamStudentStatusInfo
	pendingManualGradingCount := 0 // 统计待人工批阅的试卷数量

	for _, attempt := range attempts {
		student, err := l.svcCtx.UsersModel.FindOne(l.ctx, attempt.StudentId)
		if err != nil {
			logx.Errorf("查询学生信息失败: %v", err)
			continue
		}

		// 统计需要人工批阅的试卷数量
		if exam.RequiresManualGrading && attempt.Status == "submitted" {
			pendingManualGradingCount++
		}

		// 格式化时间，若SubmitTime为零值，则显示空字符串
		submitTimeStr := ""
		if !attempt.SubmitTime.IsZero() {
			submitTimeStr = attempt.SubmitTime.In(loc).Format("2006-01-02 15:04")
		}
		studentResults = append(studentResults, types.ExamStudentStatusInfo{
			StudentId:   attempt.StudentId,
			UserName:    student.Username,
			StudentName: student.Name,
			StartTime:   attempt.StartTime.In(loc).Format("2006-01-02 15:04"),
			SubmitTime:  submitTimeStr,
			Status:      attempt.Status,
			Score:       int(attempt.Score),
		})
	}

	// 组装返回数据
	resp := &types.ExamDetailResponse{
		ExamId:                    int(exam.Id),
		Name:                      exam.Name,
		TotalScore:                int(exam.TotalScore),
		RequiresManualGrading:     exam.RequiresManualGrading,
		StartTime:                 startTimeStr,
		EndTime:                   endTimeStr,
		ClassName:                 className,
		CreateBy:                  createBy,
		Students:                  studentResults,
		PendingManualGradingCount: pendingManualGradingCount,
	}

	return resp, nil
}
