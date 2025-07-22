package logic

import (
	"context"
	"errors"
	"exam-system/JWT"
	"exam-system/users/internal/svc"
	"exam-system/users/internal/types"
	"exam-system/users/model"
	"fmt"
	"time"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Create 创建用户

func (l *CreateLogic) Create(req *types.CreateRequest, tokenString string) (*types.CreateResponse, error) {
	// 解析JWT获取userId和userType
	_, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 检查是否为管理员
	if userType != 0 { // 0 代表管理员
		return nil, errors.New("无权限创建用户")
	}

	// 遍历用户请求数据，依次创建
	for _, user := range req.Users {
		// 生成盐
		salt := GenerateSalt()

		// 使用MD5进行加密
		hashedPasswd := HashPassword(user.Password, salt)

		// 创建用户
		_, err := l.svcCtx.UsersModel.Insert(l.ctx, &model.Users{
			Username:   user.Username, // 前端传递的username字段
			Name:       user.Name,     // 前端传递的name字段
			Passwd:     hashedPasswd,
			Salt:       salt,             // 存储加密后的密码
			Type:       int64(user.Type), // 用户类型
			CreateTime: time.Now(),       // 当前时间作为创建时间
			IsDelete:   0,                // 初始设置为0表示未删除
		})
		if err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
	}

	// 返回成功响应
	return &types.CreateResponse{
		Message: "用户创建成功",
	}, nil
}
