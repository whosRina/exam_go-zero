package logic

import (
	"context"
	"errors"
	"exam-system/JWT"
	"exam-system/users/internal/svc"
	"exam-system/users/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type LoginLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
	ctx    context.Context
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	// 查询用户
	user, err := l.svcCtx.UsersModel.FindByUsername(l.ctx, req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	hashedPassword := HashPassword(req.Password, user.Salt)
	if hashedPassword != user.Passwd {
		return nil, errors.New("密码错误")
	}

	// 生成JWT Token
	accessToken, err := jwtutil.GenerateToken(int(user.Id), int(user.Type), l.svcCtx.Config.Auth.AccessSecret, 6*time.Hour)
	if err != nil {
		return nil, err
	}

	// 返回响应
	return &types.LoginResp{
		Name:  user.Name,
		Type:  int(user.Type),
		Token: accessToken,
	}, nil
}
