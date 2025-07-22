package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/classes/internal/svc"
	"exam-system/classes/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type ClassListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClassListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClassListLogic {
	return &ClassListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *ClassListLogic) ClassList(tokenString string) (*types.ClassListResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, err
	}
	var classList []*types.ClassInfo
	switch userType {
	case 1: // 教师
		classList, err = l.svcCtx.ClassModel.FindAllByTeacherId(l.ctx, userId)
		if err != nil {
			return nil, errors.New("查询教师班级列表失败")
		}

	case 2: // 学生
		classList, err = l.svcCtx.ClassModel.FindClassesByUserId(l.ctx, userId)
		if err != nil {
			return nil, errors.New("查询学生班级列表失败")
		}

	default:
		return nil, errors.New("无权限查看班级列表")
	}

	// 构建返回数据
	var classListResponse []types.ClassInfo
	for _, class := range classList {
		classListResponse = append(classListResponse, types.ClassInfo{
			ID:          class.ID,          //班级ID
			ClassName:   class.ClassName,   // 班级名称
			ClassCode:   class.ClassCode,   // 班级邀请码
			IsJoinable:  class.IsJoinable,  // 是否允许加入
			TeacherName: class.TeacherName, // 教师名称
		})
	}

	// 返回封装好的响应数据
	return &types.ClassListResponse{
		ClassList: classListResponse,
	}, nil
}
