package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/classes/internal/svc"
	"exam-system/classes/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteClassRequest, tokenString string) (*types.DeleteClassResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 确保用户是教师才可删除班级
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限删除班级")
	}
	//
	err = l.svcCtx.ClassModel.Delete(l.ctx, req.ClassId, userId)
	if err != nil {
		return nil, errors.New("删除操作失败")
	}

	return &types.DeleteClassResponse{
		Message: "班级删除成功",
	}, nil

}
