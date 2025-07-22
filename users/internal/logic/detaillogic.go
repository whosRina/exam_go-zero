package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/users/internal/svc"
	"exam-system/users/internal/types"
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

func (l *DetailLogic) Detail(tokenString string) (*types.UserDetailResponse, error) {
	// 解析JWT获取userId和userType
	userId, _, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的 JWT")
	}

	// 查询用户信息
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 只返回用户名、姓名和用户类型
	return &types.UserDetailResponse{
		Username: user.Username,
		Name:     user.Name,
		UserType: user.Type,
	}, nil
}
