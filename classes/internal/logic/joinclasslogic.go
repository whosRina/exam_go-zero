package logic

import (
	"context"
	"errors"
	"exam-system/JWT"
	"exam-system/classes/model"

	"exam-system/classes/internal/svc"
	"exam-system/classes/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinClassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinClassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinClassLogic {
	return &JoinClassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinClassLogic) JoinClass(req *types.JoinClassRequest, tokenString string) (*types.JoinClassResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 仅允许学生加入班级
	if userType != 2 {
		return nil, errors.New("无权限加入班级")
	}

	// 查询班级是否存在
	class, err := l.svcCtx.ClassModel.FindByClassCode(l.ctx, req.ClassCode)
	if err != nil || class == nil {
		return nil, errors.New("班级不存在")
	}

	// 检查班级是否允许加入
	if class.IsJoinable == 0 {
		return nil, errors.New("该班级不允许加入")
	}

	// 检查用户是否已经在班级中
	member, err := l.svcCtx.ClassMemberModel.FindByUserAndClass(l.ctx, int(userId), int(class.Id))
	if err == nil && member != nil {
		return nil, errors.New("你已经在该班级中")
	}

	// 插入班级成员表
	_, err = l.svcCtx.ClassMemberModel.Insert(l.ctx, &model.ClassMember{
		ClassId: class.Id,
		UserId:  userId,
	})
	if err != nil {
		logx.Error("加入班级失败:", err)
		return nil, errors.New("加入班级失败")
	}

	return &types.JoinClassResponse{
		Message: "成功加入班级",
	}, nil
}
