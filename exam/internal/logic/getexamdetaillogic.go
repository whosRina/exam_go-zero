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

type GetExamDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExamDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExamDetailLogic {
	return &GetExamDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExamDetailLogic) GetExamDetail(req *types.GetExamDetailRequest, tokenString string) (*types.GetExamDetailResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 仅允许学生访问
	if userType != 2 {
		return nil, errors.New("无权限通过该接口查看该考试详情")
	}

	// 查询考试信息
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, req.ExamId)
	if err != nil {
		return nil, errors.New("考试不存在")
	}

	// 验证学生是否属于考试班级
	isMember, err := l.svcCtx.ClassMemberModel.Exists(l.ctx, exam.ClassId, userId)
	if err != nil {
		return nil, errors.New("无法验证班级成员")
	}
	if !isMember {
		return nil, errors.New("你不属于此考试班级，无法查看详情")
	}

	// 查询班级名称
	className := ""
	class, err := l.svcCtx.ClassModel.FindOne(l.ctx, exam.ClassId)
	if err == nil {
		className = class.Name
	} else {
		logx.Errorf("查询班级名称失败: %v", err)
	}

	// 查询创建者名称
	createBy := ""
	teacher, err := l.svcCtx.UsersModel.FindOne(l.ctx, exam.TeacherId)
	if err == nil {
		createBy = teacher.Name
	} else {
		logx.Errorf("查询创建者名称失败: %v", err)
	}

	// 格式化时间，使用 "2006-01-02 15:04" 格式（上海时区）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTimeStr := exam.StartTime.In(loc).Format("2006-01-02 15:04")
	endTimeStr := exam.EndTime.In(loc).Format("2006-01-02 15:04")

	attempt, err := l.svcCtx.ExamAttemptModel.FindAttemptByExamAndStudent(l.ctx, exam.Id, userId)

	var score int // 默认分数为0
	if err != nil {
		logx.Errorf("查询考试尝试记录失败: %v", err)
	} else if attempt != nil {
		score = int(attempt.Score)
	}
	var status string
	if attempt != nil {
		status = attempt.Status
	} else {
		status = "not_started"
	}

	// 组装返回数据
	resp := &types.GetExamDetailResponse{
		ExamId:                int(exam.Id),
		Name:                  exam.Name,
		TotalScore:            int(exam.TotalScore),
		RequiresManualGrading: exam.RequiresManualGrading,
		StartTime:             startTimeStr,
		EndTime:               endTimeStr,
		CanViewResults:        exam.CanViewResults,
		ClassName:             className,
		CreateBy:              createBy,
		Score:                 score,
		Status:                status,
	}

	return resp, nil
}
