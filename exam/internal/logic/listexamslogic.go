package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"exam-system/exam/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListExamsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListExamsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListExamsLogic {
	return &ListExamsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListExamsLogic) ListExams(tokenString string) (*types.ExamListResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, err
	}

	var exams []*model.Exam
	switch userType {
	case 1: // 教师
		// 查询教师发布的所有考试
		exams, err = l.svcCtx.ExamModel.FindByTeacherId(l.ctx, userId)
		if err != nil {
			logx.Errorf("查询教师考试列表失败: %v", err)
			return nil, errors.New("获取考试列表失败")
		}

	case 2: // 学生
		// 先获取学生加入的班级列表
		classIds, err := l.svcCtx.ClassMemberModel.FindByUser(l.ctx, userId)
		if err != nil {
			logx.Errorf("查询学生班级列表失败: %v", err)
			return nil, errors.New("获取班级信息失败")
		}

		// 根据班级ID查询关联的考试
		if len(classIds) > 0 {
			exams, err = l.svcCtx.ExamModel.FindByClassIds(l.ctx, classIds)
			if err != nil {
				logx.Errorf("查询班级考试失败: %v", err)
				return nil, errors.New("获取考试列表失败")
			}
		}

	default:
		return nil, errors.New("无权限访问")
	}

	// 转换时间格式并返回
	response := make([]types.ExamInfo, 0, len(exams))
	for _, exam := range exams {
		// 查询班级名称
		class, err := l.svcCtx.ClassModel.FindOne(l.ctx, exam.ClassId)
		className := ""
		if err == nil {
			className = class.Name
		} else {
			logx.Errorf("查询班级名称失败: %v", err)
		}

		// 查询创建者名称
		teacher, err := l.svcCtx.UsersModel.FindOne(l.ctx, exam.TeacherId)
		creatorName := ""
		if err == nil {
			creatorName = teacher.Name
		} else {
			logx.Errorf("查询创建者名称失败: %v", err)
		}
		response = append(response, types.ExamInfo{
			Id:                    int(exam.Id),
			Name:                  exam.Name,
			ClassId:               int(exam.ClassId),
			ExamType:              exam.ExamType,
			TotalScore:            int(exam.TotalScore),
			RequiresManualGrading: exam.RequiresManualGrading,
			StartTime:             exam.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:               exam.EndTime.Format("2006-01-02 15:04:05"),
			CanViewResults:        exam.CanViewResults,
			PaperId:               exam.PaperId,
			PaperRuleId:           exam.PaperRuleId,
			ClassName:             className,
			CreateBy:              creatorName,
		})
	}

	return &types.ExamListResponse{
		ExamList: response,
	}, nil
}
