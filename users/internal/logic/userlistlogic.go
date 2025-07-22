package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/users/internal/svc"
	"exam-system/users/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(tokenString string) (*types.UserListResponse, error) {
	// 解析JWT获取userId和userType
	_, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 检查是否为管理员
	if userType != 0 { // 0代表管理员
		return nil, errors.New("无权限创建用户")
	}
	// 查询所有未删除的用户
	users, err := l.svcCtx.UsersModel.FindAll(l.ctx)
	if err != nil {
		return nil, errors.New("查询用户列表失败")
	}

	// 构造返回数据
	var userList []types.UserInfo
	for _, user := range users {
		userList = append(userList, types.UserInfo{
			Id:       user.Id,
			Name:     user.Name,
			Username: user.Username,
			Type:     int(user.Type),
		})
	}

	// 返回响应
	return &types.UserListResponse{
		UserList: userList,
	}, nil
}
