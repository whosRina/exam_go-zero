package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/users/internal/types"
	"fmt"

	"exam-system/users/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateUserRequest, tokenString string) (*types.UpdateUserResponse, error) {
	// 解析JWT获取userId和userType
	_, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 检查是否为管理员
	if userType != 0 { // 0 代表管理员
		return nil, errors.New("无权限更新用户")
	}

	// 查询用户是否存在
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("用户不存在:%v", err)
	}

	// 更新用户基本信息
	user.Name = req.Name
	user.Username = req.Username
	user.Type = int64(req.Type)

	// 如果密码不为空，则更新密码
	if req.Password != "" {
		salt := GenerateSalt() // 生成新的盐值
		hashedPasswd := HashPassword(req.Password, salt)
		user.Passwd = hashedPasswd
		user.Salt = salt
	}

	// 更新数据库
	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("更新用户失败:%v", err)
	}

	// 返回成功信息
	return &types.UpdateUserResponse{
		Message: "用户信息更新成功",
	}, nil
}
