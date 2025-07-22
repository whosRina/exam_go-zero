package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/classes/internal/svc"
	"exam-system/classes/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMemberLogic {
	return &DeleteMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMemberLogic) DeleteMember(req *types.DeleteClassMemberRequest, tokenString string) (*types.DeleteClassMemberResponse, error) {

	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 确保用户是教师才可删除成员
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限删除该成员")
	}
	// 验证该成员是否属于教师的班级
	isMember, err := l.svcCtx.ClassMemberModel.IsMemberInTeacherClass(l.ctx, int(userId), int(req.MemberId))
	if err != nil {
		return nil, errors.New("数据库查询失败")
	}
	if !isMember {
		return nil, errors.New("该成员不属于您的班级，无法操作")
	}

	// 执行删除操作
	err = l.svcCtx.ClassMemberModel.Delete(l.ctx, req.MemberId)
	if err != nil {
		return nil, errors.New("删除成员失败")
	}

	// 返回成功响应
	return &types.DeleteClassMemberResponse{
		Message: "删除成功",
	}, nil

}
