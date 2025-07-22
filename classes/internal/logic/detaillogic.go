package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/classes/internal/types"

	"exam-system/classes/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.ClassDetailRequest, tokenString string) (*types.ClassDetailResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 根据用户类型鉴别用户是否有权限访问班级详情
	if userType == 1 { // 教师
		// 教师验证：检查班级是否存在且该教师是否是该班级的创建者
		class, err := l.svcCtx.ClassModel.FindOne(l.ctx, req.ClassId)
		if err != nil {
			return nil, errors.New("查询班级详情失败")
		}
		// 确保查询到的班级属于该教师
		if class.TeacherId != userId {
			return nil, errors.New("无权限查看此班级")
		}
	} else if userType == 2 { // 学生
		// 检查用户是否在该班级的成员表中
		member, err := l.svcCtx.ClassMemberModel.FindByUserAndClass(l.ctx, int(userId), int(req.ClassId))
		if err != nil {
			return nil, errors.New("查询班级成员失败")
		}
		if member == nil {
			return nil, errors.New("你不是该班级的成员")
		}
	}
	// 获取班级详情
	classDetail, err := l.svcCtx.ClassModel.FindClassDetailByClassID(l.ctx, req.ClassId)
	if err != nil {
		logx.Error("获取班级详情失败:", err)
		return nil, errors.New("查询班级详情失败")
	}
	// 构建返回数据
	classDetailResponse := &types.ClassDetailResponse{
		ClassInfo:    classDetail.ClassInfo,
		ClassMembers: classDetail.ClassMembers,
	}

	return classDetailResponse, nil
}
